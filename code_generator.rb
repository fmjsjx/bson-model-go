require 'set'
require 'yaml'


def to_underscore(name)
  field_name = name.gsub('/[A-Z]/') do |match| "_#{match}" end
  field_name.downcase
end

def to_camel(name)
  name[0].upcase << name[1..-1]
end

def to_small_camel(name)
  name[0].downcase << name[1..-1]
end

def tabs(n, value = '')
  "\t" * n + value + "\n"
end

def fix_space(value, size)
  if value.size < size
    value + ' ' * (size - value.size)
  else
    value
  end
end

def map_type(key_type)
  case key_type
  when 'int'
    'bsonmodel.IntObjectMapModel'
  when 'string'
    'bsonmodel.StringObjectMapModel'
  else
    raise "supported key type `#{key_type}` for map"
  end
end

def simple_map_type(key_type)
  case key_type
  when 'int'
    'bsonmodel.IntSimpleMapModel'
  when 'string'
    'bsonmodel.StringSimpleMapModel'
  else
    raise "supported key type `#{key_type}` for simple map"
  end
end

def map_factory(key_type)
  case key_type
  when 'int'
    'bsonmodel.NewIntObjectMapModel'
  when 'string'
    'bsonmodel.NewStringObjectMapModel'
  else
    raise "supported key type `#{key_type}` for map"
  end
end

def map_value_factory(key_type)
  case key_type
  when 'int'
    'bsonmodel.IntObjectMapValueFactory'
  when 'string'
    'bsonmodel.StringObjectMapValueFactory'
  else
    raise "supported key type `#{key_type}` for map"
  end
end

def simple_map_factory(key_type)
  case key_type
  when 'int'
    'bsonmodel.NewIntSimpleMapModel'
  when 'string'
    'bsonmodel.NewStringSimpleMapModel'
  else
    raise "supported key type `#{key_type}` for simple map"
  end
end

def simple_value_type(value_type)
  case value_type
  when 'int'
    'bsonmodel.IntValueType()'
  when 'string'
    'bsonmodel.StringValeType()'
  when 'float64'
    'bsonmodel.Float64ValueType()'
  when 'bool'
    'bsonmodel.BoolValueType()'
  when 'datetime'
    'bsonmodel.DateTimeSimpleValueType()'
  when 'date'
    'bsonmodel.DateSimpleValueType()'
  else
    raise "unsupported value type `#{value_type}` for simple map"
  end
end

def map_value_type(key_type)
  case key_type
  when 'int'
    'bsonmodel.IntObjectMapValueModel'
  when 'string'
    'bsonmodel.StringObjectMapValueModel'
  else
    raise "supported key type `#{key_type}` for map value"
  end
end

def map_value_struct(key_type)
  case key_type
  when 'int'
    'bsonmodel.BaseIntObjectMapValue'
  when 'string'
    'bsonmodel.BaseStringObjectMapValue'
  else
    raise "supported key type `#{key_type}` for map value"
  end
end

def fill_imports(code, cfg)
  stds = Set.new
  others = Set.new
  aliases = {'github.com/json-iterator/go' => 'jsoniter'}
  stds << 'unsafe'
  others << 'github.com/bits-and-blooms/bitset'
	others << 'github.com/fmjsjx/bson-model-go/bsonmodel'
	others << 'github.com/json-iterator/go'
	others << 'go.mongodb.org/mongo-driver/bson'
  if cfg['fields'].any? { |field| %w(datetime date).include? field['type'] }
    stds << 'time'
    others << 'go.mongodb.org/mongo-driver/bson/primitive'
  end
  code << "import (\n"
  stds.sort.each do |v|
    if aliases.include? v
      code << tabs(1, "#{aliases[v]} \"#{v}\"")
    else
      code << tabs(1, "\"#{v}\"")
    end
  end
  code << "\n"
  others.sort.each do |v|
    if aliases.include? v
      code << tabs(1, "#{aliases[v]} \"#{v}\"")
    else
      code << tabs(1, "\"#{v}\"")
    end
  end
  code << ")\n\n"
end

def fill_interface(code, super_interface, cfg)
  code << "type #{cfg['name']} interface {\n"
  code << tabs(1, "#{super_interface}")
  fields = cfg['fields']
  fields.each do |field|
    name = field['name']
    camel = to_camel(name)
    case field['type']
    when 'int'
      code << tabs(1, "#{camel}() int")
      unless field['virtual'] == true
        code << tabs(1, "Set#{camel}(#{name} int)")
        if field['increase'] == true
          code << tabs(1, "Increase#{camel}() int")
        end
        if field['add'] == true
          code << tabs(1, "Add#{camel}(#{name} int) int")
        end
      end
    when 'string'
      code << tabs(1, "#{camel}() string")
      unless field['virtual'] == true
        code << tabs(1, "Set#{camel}(#{name} string)")
      end
    when 'float64'
      code << tabs(1, "#{camel}() float64")
      unless field['virtual'] == true
        code << tabs(1, "Set#{camel}(#{name} float64)")
      end
    when 'datetime'
      code << tabs(1, "#{camel}() time.Time")
      unless field['virtual'] == true
        code << tabs(1, "Set#{camel}(#{name} time.Time)")
      end
    when 'date'
      code << tabs(1, "#{camel}() time.Time")
      unless field['virtual'] == true
        code << tabs(1, "Set#{camel}(#{name} time.Time)")
        code << tabs(1, "Set#{camel}Number(#{name} int)")
      end
    when 'object'
      code << tabs(1, "#{camel}() #{field['model']}")
    when 'map'
      key_type = field['key']
      value_type = field['value']
      code << tabs(1, "#{camel}() #{map_type(key_type)}")
      if field.has_key? 'quick-access-method'
        code << tabs(1, "#{field['quick-access-method']}(id #{key_type}) #{value_type}")
      elsif camel.end_with? 's'
        code << tabs(1, "#{camel[0..-2]}(id #{key_type}) #{value_type}")
      end
    when 'simple-map'
      key_type = field['key']
      code << tabs(1, "#{camel}() #{simple_map_type(key_type)}")
    else
      raise "unsupported field type `#{field['type']}` on #{cfg['name']}.#{field['name']}"
    end
  end
  code << "}\n\n"
end

def fill_const(code, cfg)
  fields = cfg['fields']
  consts = []
  max_len = 0
  fields.each do |field|
    next if field['virtual'] == true
    const_name = "Bname#{cfg['name']}#{to_camel(field['name'])}"
    consts << [const_name, field['bname']]
    if const_name.size > max_len
      max_len = const_name.size
    end
  end
  unless consts.empty?
    code << "const (\n"
    consts.each do |const|
      code << tabs(1, "#{fix_space(const[0], max_len)} = \"#{const[1]}\"")
    end
    code << ")\n\n"
  end
end

def all_object(cfg)
  parent = cfg['parent']
  if parent.nil?
    return false
  end
  if parent['type'] == 'root'
    return true
  else
    return all_object(parent)
  end
end

def fill_struct(code, cfg, super_struct=nil)
  fields = cfg['fields']
  max_len = "updatedFields".size
  fields.each do |field|
    next if field['virtual'] == true
    len = field['name'].size
    if len > max_len
      max_len = len
    end
  end
  code << "type default#{cfg['name']} struct {\n"
  unless super_struct.nil?
    code << tabs(1, super_struct)
  end
  code << tabs(1, "#{fix_space('updatedFields', max_len)} *bitset.BitSet")
  if cfg['type'] == 'object'
    parent = cfg['parent']
    code << tabs(1, "#{fix_space('parent', max_len)} #{parent['name']}")
  end
  fields.each do |field|
    next if field['virtual'] == true
    name = field['name']
    case field['type']
    when 'int'
      code << tabs(1, "#{fix_space(name, max_len)} int")
    when 'string'
      code << tabs(1, "#{fix_space(name, max_len)} string")
    when 'float64'
      code << tabs(1, "#{fix_space(name, max_len)} float64")
    when 'datetime'
      code << tabs(1, "#{fix_space(name, max_len)} time.Time")
    when 'time'
      code << tabs(1, "#{fix_space(name, max_len)} time.Time")
    when 'object'
      code << tabs(1, "#{fix_space(name, max_len)} #{field['model']}")
    when 'map'
      key_type = field['key']
      code << tabs(1, "#{fix_space(name, max_len)} #{map_type(key_type)}")
    when 'simple-map'
      key_type = field['key']
      code << tabs(1, "#{fix_space(name, max_len)} #{simple_map_type(key_type)}")
    else
      raise "unsupported field type `#{field['type']}` on #{cfg['name']}.#{name}"
    end
  end
  code << "}\n\n"
end

def fill_to_bson(code, cfg)
  code << "func (self *default#{cfg['name']}) ToBson() interface{} {\n"
  code << tabs(1, "return self.ToDocument()")
  code << "}\n\n"
end

def fill_to_data(code, cfg)
  code << "func (self *default#{cfg['name']}) ToData() interface{} {\n"
  code << tabs(1, "data := make(map[string]interface{})")
  cfg['fields'].each do |field|
    next if field['virtual'] == true
    type = field['type']
    case 
    when %w(object map simple-map).include?(type)
      code << tabs(1, "data[\"#{field['bname']}\"] = self.#{field['name']}.ToData()")
    when type == 'datetime'
      code << tabs(1, "data[\"#{field['bname']}\"] = self.#{field['name']}.UnixMilli()")
    when type == 'date'
      code << tabs(1, "data[\"#{field['bname']}\"] = bsonmodel.DateToNumber(self.#{field['name']})")
    else
      code << tabs(1, "data[\"#{field['bname']}\"] = self.#{field['name']}")
    end
  end
  code << tabs(1, "return data")
  code << "}\n\n"
end

def fill_load_jsoniter(code, cfg, is_root = false)
  code << "func (self *default#{cfg['name']}) LoadJsoniter(any jsoniter.Any) error {\n"
  code << tabs(1, "if any.ValueType() != jsoniter.ObjectValue {")
  if is_root
    code << tabs(2, "self.Reset()")
  end
  code << tabs(2, "return nil")
  code << tabs(1, "}")
  cfg['fields'].each do |field|
    next if field['virtual'] == true
    name = field['name']
    bname = field['bname']
    case field['type']
    when 'int'
      default = field.has_key?('default') ? field['default'].to_i : 0
      code << tabs(1, "#{name}, err := bsonmodel.AnyIntValue(any.Get(\"#{bname}\"), #{default})")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "self.#{name} = #{name}")
    when 'string'
      default = field.has_key?('default') ? field['default'].to_s : ''
      code << tabs(1, "#{name}, err := bsonmodel.AnyStringValue(any.Get(\"#{bname}\"), \"#{default}\")")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "self.#{name} = #{name}")
    when 'float64'
      default = field.has_key?('default') ? field['default'] : '0'
      code << tabs(1, "#{name}, err := bsonmodel.AnyFloat64Value(any.Get(\"#{bname}\"), #{default})")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "self.#{name} = #{name}")
    when 'datetime'
      code << tabs(1, "#{name}, err := bsonmodel.AnyDateTimeValue(any.Get(\"#{bname}\"))")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "self.#{name} = #{name}")
    when 'date'
      code << tabs(1, "#{name}, err := bsonmodel.AnyDateValue(any.Get(\"#{bname}\"))")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "self.#{name} = #{name}")
    when 'object'
      code << tabs(1, "#{name} := any.Get(\"#{bname}\")")
      code << tabs(1, "if #{name}.ValueType() == jsoniter.ObjectValue {")
      code << tabs(2, "err = self.#{name}.LoadJsoniter(#{name})")
      code << tabs(2, "if err != nil {")
      code << tabs(3, "return err")
      code << tabs(2, "}")
      code << tabs(1, "}")
    when 'map'
      code << tabs(1, "#{name} := any.Get(\"#{bname}\")")
      code << tabs(1, "if #{name}.ValueType() == jsoniter.ObjectValue {")
      code << tabs(2, "err = self.#{name}.LoadJsoniter(#{name})")
      code << tabs(2, "if err != nil {")
      code << tabs(3, "return err")
      code << tabs(2, "}")
      code << tabs(1, "} else {")
      code << tabs(2, "self.#{name}.Clear()")
      code << tabs(1, "}")
    when 'simple-map'
      code << tabs(1, "#{name} := any.Get(\"#{bname}\")")
      code << tabs(1, "if #{name}.ValueType() == jsoniter.ObjectValue {")
      code << tabs(2, "err = self.#{name}.LoadJsoniter(#{name})")
      code << tabs(2, "if err != nil {")
      code << tabs(3, "return err")
      code << tabs(2, "}")
      code << tabs(1, "} else {")
      code << tabs(2, "self.#{name}.Reset()")
      code << tabs(1, "}")
    end
  end
  if is_root
    code << tabs(1, "self.Reset()")
  end
  code << tabs(1, "return nil")
  code << "}\n\n"
end

def fill_reset(code, cfg)
  code << "func (self *default#{cfg['name']}) Reset() {\n"
  cfg['fields'].each do |field|
    next if field['virtual'] == true
    if %w(object map simple-map).include? field['type']
      code << tabs(1, "self.#{field['name']}.Reset()")
    end
  end
  code << tabs(1, "self.updatedFields.ClearAll()")
  code << "}\n\n"
end

def fill_any_updated(code, cfg)
  code << "func (self *default#{cfg['name']}) AnyUpdated() bool {\n"
  any_updateds = []
  cfg['fields'].each do |field|
    next if field['virtual'] == true
    if %w(object map simple-map).include? field['type']
      any_updateds << "self.#{field['name']}.AnyUpdated()"
    end
  end
  if any_updateds.empty?
    code << tabs(1, "return self.updatedFields.Any()")
  else
    code << tabs(1, "return self.updatedFields.Any() || #{any_updateds.join(' || ')}")
  end
  code << "}\n\n"
end

def fill_any_deleted(code, cfg)
  code << "func (self *default#{cfg['name']}) AnyDeleted() bool {\n"
  code << tabs(1, "return self.DeletedSize() > 0")
  code << "}\n\n"
end

def fill_append_updates(code, cfg)
  code << "func (self *default#{cfg['name']}) AppendUpdates(updates bson.M) bson.M {\n"
  code << tabs(1, "dset := bsonmodel.FixedEmbedded(updates, \"$set\")")
  if cfg['type'] == 'root'
    code << tabs(1, "updatedFields := self.updatedFields")
    cfg['fields'].each_with_index do |field, index|
      next if field['virtual'] == true
      name = field['name']
      bname = field['bname']
      if %w(object map simple-map).include? field['type']
        code << tabs(1, "if self.#{name}.AnyUpdated() {")
        code << tabs(2, "self.#{name}.AppendUpdates(updates)")
      else
        code << tabs(1, "if updatedFields.Test(#{index + 1}) {")
        case field['type']
        when 'datetime'
          code << tabs(2, "dset[\"#{bname}\"] = primitive.NewDateTimeFromTime(self.#{name})")
        when 'date'
          code << tabs(2, "dset[\"#{bname}\"] = bsonmodel.DateToNumber(self.#{name})")
        else
          code << tabs(2, "dset[\"#{bname}\"] = self.#{name}")
        end
      end
      code << tabs(1, "}")
    end
  else
    code << tabs(1, "xpath := self.XPath()")
    code << tabs(1, "if self.FullyUpdate() {")
    code << tabs(2, "dset[xpath.Value()] = self.ToDocument()")
    code << tabs(1, "} else {")
    code << tabs(2, "updatedFields := self.updatedFields")
    cfg['fields'].each_with_index do |field, index|
      next if field['virtual'] == true
      name = field['name']
      bname = field['bname']
      if %w(object map simple-map).include? field['type']
        code << tabs(2, "if self.#{name}.AnyUpdated() {")
        code << tabs(3, "self.#{name}.AppendUpdates(updates)")
      else
        code << tabs(2, "if updatedFields.Test(#{index + 1}) {")
        case field['type']
        when 'datetime'
          code << tabs(3, "dset[xpath.Resolve(\"#{bname}\").Value()] = primitive.NewDateTimeFromTime(self.#{name}}")
        when 'date'
          code << tabs(3, "dset[xpath.Resolve(\"#{bname}\").Value()] = bsonmodel.DateToNumber(self.#{name}}")
        else
          code << tabs(3, "dset[xpath.Resolve(\"#{bname}\").Value()] = self.#{name}")
        end
      end
      code << tabs(2, "}")
    end
    code << tabs(1, "}")
  end
  code << tabs(1, "return updates")
  code << "}\n\n"
end

def fill_to_document(code, cfg)
  code << "func (self *default#{cfg['name']}) ToDocument() bson.M {\n"
  code << tabs(1, "doc := bson.M{}")
  cfg['fields'].each do |field|
    next if field['virtual'] == true
    name = field['name']
    bname = field['bname']
    if %w(object map simple-map).include? field['type']
      code << tabs(1, "doc[\"#{bname}\"] = self.#{name}.ToBson()")
    else
      case field['type']
      when 'datetime'
        code << tabs(1, "doc[\"#{bname}\"] = primitive.NewDateTimeFromTime(self.#{name})")
      when 'date'
        code << tabs(1, "doc[\"#{bname}\"] = bsonmodel.DateToNumber(self.#{name})")
      else
        code << tabs(1, "doc[\"#{bname}\"] = self.#{name}")
      end
    end
  end
  code << tabs(1, "return doc")
  code << "}\n\n"
end

def fill_load_document(code, cfg, is_root = false)
  code << "func (self *default#{cfg['name']}) LoadDocument(document bson.M) error {\n"
  cfg['fields'].each do |field|
    next if field['virtual'] == true
    name = field['name']
    bname = field['bname']
    case field['type']
    when 'int'
      default = field.has_key?('default') ? field['default'].to_i : 0
      code << tabs(1, "#{name}, err := bsonmodel.IntValue(document, \"#{bname}\", #{default})")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "self.#{name} = #{name}")
    when 'string'
      default = field.has_key?('default') ? field['default'].to_s : ''
      code << tabs(1, "#{name}, err := bsonmodel.StringValue(document, \"#{bname}\", \"#{default}\")")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "self.#{name} = #{name}")
    when 'float64'
      default = field.has_key?('default') ? field['default'] : '0'
      code << tabs(1, "#{name}, err := bsonmodel.Float64Value(document, \"#{bname}\", #{default})")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "self.#{name} = #{name}")
    when 'datetime'
      code << tabs(1, "#{name}, err := bsonmodel.DateTimeValue(document, \"#{bname}\")")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "self.#{name} = #{name}")
    when 'date'
      code << tabs(1, "#{name}, err := bsonmodel.DateValue(document, \"#{bname}\")")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "self.#{name} = #{name}")
    when 'object'
      code << tabs(1, "#{name}, err := bsonmodel.EmbeddedValue(document, \"#{bname}\")")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "if #{name} != nil {")
      code << tabs(2, "err = self.#{name}.LoadDocument(#{name})")
      code << tabs(2, "if err != nil {")
      code << tabs(3, "return err")
      code << tabs(2, "}")
      code << tabs(1, "}")
    when 'map'
      code << tabs(1, "#{name}, err := bsonmodel.EmbeddedValue(document, \"#{bname}\")")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "if #{name} != nil {")
      code << tabs(2, "err = self.#{name}.LoadDocument(#{name})")
      code << tabs(2, "if err != nil {")
      code << tabs(3, "return err")
      code << tabs(2, "}")
      code << tabs(1, "} else {")
      code << tabs(2, "self.#{name}.Clear()")
      code << tabs(1, "}")
    when 'simple-map'
      code << tabs(1, "#{name}, err := bsonmodel.EmbeddedValue(document, \"#{bname}\")")
      code << tabs(1, "if err != nil {")
      code << tabs(2, "return err")
      code << tabs(1, "}")
      code << tabs(1, "if #{name} != nil {")
      code << tabs(2, "err = self.#{name}.LoadDocument(#{name})")
      code << tabs(2, "if err != nil {")
      code << tabs(3, "return err")
      code << tabs(2, "}")
      code << tabs(1, "} else {")
      code << tabs(2, "self.#{name}.Clear()")
      code << tabs(1, "}")
    end
  end
  if is_root
    code << tabs(1, "self.Reset()")
  end
  code << tabs(1, "return nil")
  code << "}\n\n"
end

def fill_deleted_size(code, cfg)
  code << "func (self *default#{cfg['name']}) DeletedSize() int {\n"
  children = cfg['fields'].select do |field|
    %w(object map simple-map).include? field['type']
  end.map do |field|
    field['name']
  end
  if children.empty?
    code << tabs(1, "return 0")
  else
    code << tabs(1, "n := 0")
    children.each do |name|
      code << tabs(1, "if self.#{name}.AnyDeleted() {")
      code << tabs(2, "n += 1")
      code << tabs(1, "}")
    end
    code << tabs(1, "return n")
  end
  code << "}\n\n"
end

def fill_fully_update(code, cfg, is_root = false)
  code << "func (self *default#{cfg['name']}) FullyUpdate() bool {\n"
  if is_root
    code << tabs(1, "return false")
  else
    code << tabs(1, "return self.updatedFields.Test(0)")
  end
  code << "}\n\n"
  code << "func (self *default#{cfg['name']}) SetFullyUpdate(fullyUpdate bool) {\n"
  if is_root
    code << tabs(1, "// no effect")
  else
    code << tabs(1, "if fullyUpdate {")
    code << tabs(2, "self.updatedFields.Set(0)")
    code << tabs(1, "} else {")
    code << tabs(2, "self.updatedFields.DeleteAt(0)")
    code << tabs(1, "}")
  end
  code << "}\n\n"
end

def fill_to_sync(code, cfg, is_root = false)
  code << "func (self *default#{cfg['name']}) ToSync() interface{} {\n"
  unless is_root
    code << tabs(1, "if self.FullyUpdate() {")
    code << tabs(2, "return self")
    code << tabs(1, "}")
  end
  code << tabs(1, "sync := make(map[string]interface{})")
  lines = []
  cfg['fields'].each_with_index do |field, index|
    next if field['json-ignore'] == true
    name = field['name']
    if %w(object map simple-map).include? field['type']
      lines << tabs(1, "if self.#{name}.AnyUpdated() {")
      lines << tabs(2, "sync[\"#{name}\"] = self.#{name}.ToSync()")
    else
      lines << tabs(1, "if updatedFields.Test(#{index + 1}) {")
      case field['type']
      when 'datetime'
        if field['virtual'] == true
          lines << tabs(2, "sync[\"#{name}\"] = self.#{to_camel(name)}().Unix()")
        else
          lines << tabs(2, "sync[\"#{name}\"] = self.#{name}.Unix()")
        end
      when 'date'
        if field['virtual'] == true
          lines << tabs(2, "sync[\"#{name}\"] = bsonmodel.DateToNumber(self.#{to_camel(name)}())")
        else
          lines << tabs(2, "sync[\"#{name}\"] = bsonmodel.DateToNumber(self.#{name})")
        end
      else
        if field['virtual'] == true
          lines << tabs(2, "sync[\"#{name}\"] = self.#{to_camel(name)}()")
        else
          lines << tabs(2, "sync[\"#{name}\"] = self.#{name}")
        end
      end
    end
    lines << tabs(1, "}")
  end
  unless lines.empty?
    code << tabs(1, "updatedFields := self.updatedFields")
    code << lines.join
  end
  code << tabs(1, "return sync")
  code << "}\n\n"
end

def fill_to_delete(code, cfg)
  code << "func (self *default#{cfg['name']}) ToDelete() interface{} {\n"
  code << tabs(1, "delete := make(map[string]interface{})")
  cfg['fields'].each_with_index do |field, index|
    next if field['json-ignore'] == true
    name = field['name']
    if %w(object map simple-map).include? field['type']
      code << tabs(1, "if self.#{name}.AnyDeleted() {")
      code << tabs(2, "delete[\"#{name}\"] = self.#{name}.ToDelete()")
      code << tabs(1, "}")
    end
  end
  code << tabs(1, "return delete")
  code << "}\n\n"
end

def fill_xetters(code, cfg)
  fields = cfg['fields']
  fields.each_with_index do |field, index|
    name = field['name']
    camel = to_camel(name)
    case field['type']
    when 'int'
      if field['virtual'] == true
        code << "func (self *default#{cfg['name']}) #{camel}() int {\n"
        unless field.has_key? 'formula'
          raise "missing required field `formula` on #{cfg['name']}.#{name}"
        end
        code << tabs(1, "return #{field['formula']}")
        code << "}\n\n"
      else
        code << "func (self *default#{cfg['name']}) #{camel}() int {\n"
        code << tabs(1, "return self.#{name}")
        code << "}\n\n"
        code << "func (self *default#{cfg['name']}) Set#{camel}(#{name} int) {\n"
        code << tabs(1, "if self.#{name} != #{name} {")
        code << tabs(2, "self.#{name} = #{name}")
        code << tabs(2, "self.updatedFields.Set(#{index + 1})")
        code << tabs(1, "}")
        code << "}\n\n"
        if field['increase'] == true
          code << "func (self *default#{cfg['name']}) Increase#{camel}() int {\n"
          code << tabs(1, "#{name} := self.#{name} + 1")
          code << tabs(1, "self.#{name} = #{name}")
          code << tabs(1, "self.updatedFields.Set(#{index + 1})")
          code << tabs(1, "return #{name}")
          code << "}\n\n"
        end
        if field['add'] == true
          code << tabs(1, "Add#{camel}(#{name} int) int")
          code << "func (self *default#{cfg['name']}) Add#{camel}(#{name} int) int {\n"
          code << tabs(1, "new_#{name} := self.#{name} + #{name}")
          code << tabs(1, "self.#{name} = new_#{name}")
          code << tabs(1, "self.updatedFields.Set(#{index + 1})")
          code << tabs(1, "return new_#{name}")
          code << "}\n\n"
        end
      end
    when 'string'
      if field['virtual'] == true
        code << "func (self *default#{cfg['name']}) #{camel}() string {\n"
        unless field.has_key? 'formula'
          raise "missing required field `formula` on #{cfg['name']}.#{name}"
        end
        code << tabs(1, "return #{field['formula']}")
        code << "}\n\n"
      else
        code << "func (self *default#{cfg['name']}) #{camel}() string {\n"
        code << tabs(1, "return self.#{name}")
        code << "}\n\n"
        code << "func (self *default#{cfg['name']}) Set#{camel}(#{name} string) {\n"
        code << tabs(1, "if self.#{name} != #{name} {")
        code << tabs(2, "self.#{name} = #{name}")
        code << tabs(2, "self.updatedFields.Set(#{index + 1})")
        code << tabs(1, "}")
        code << "}\n\n"
      end
    when 'float64'
      if field['virtual'] == true
        code << "func (self *default#{cfg['name']}) #{camel}() float64 {\n"
        unless field.has_key? 'formula'
          raise "missing required field `formula` on #{cfg['name']}.#{name}"
        end
        code << tabs(1, "return #{field['formula']}")
        code << "}\n\n"
      else
        code << "func (self *default#{cfg['name']}) #{camel}() float64 {\n"
        code << tabs(1, "return self.#{name}")
        code << "}\n\n"
        code << "func (self *default#{cfg['name']}) Set#{camel}(#{name} float64) {\n"
        code << tabs(1, "if self.#{name} != #{name} {")
        code << tabs(2, "self.#{name} = #{name}")
        code << tabs(2, "self.updatedFields.Set(#{index + 1})")
        code << tabs(1, "}")
        code << "}\n\n"
      end
    when 'datetime'
      if field['virtual'] == true
        code << "func (self *default#{cfg['name']}) #{camel}() time.Time {\n"
        unless field.has_key? 'formula'
          raise "missing required field `formula` on #{cfg['name']}.#{name}"
        end
        code << tabs(1, "return #{field['formula']}")
        code << "}\n\n"
      else
        code << "func (self *default#{cfg['name']}) #{camel}() time.Time {\n"
        code << tabs(1, "return self.#{name}")
        code << "}\n\n"
        code << "func (self *default#{cfg['name']}) Set#{camel}(#{name} time.Time) {\n"
        code << tabs(1, "if self.#{name} != #{name} {")
        code << tabs(2, "self.#{name} = #{name}")
        code << tabs(2, "self.updatedFields.Set(#{index + 1})")
        code << tabs(1, "}")
        code << "}\n\n"
      end
    when 'date'
      if field['virtual'] == true
        code << "func (self *default#{cfg['name']}) #{camel}() time.Time {\n"
        unless field.has_key? 'formula'
          raise "missing required field `formula` on #{cfg['name']}.#{name}"
        end
        code << tabs(1, "return #{field['formula']}")
        code << "}\n\n"
      else
        code << "func (self *default#{cfg['name']}) #{camel}() time.Time {\n"
        code << tabs(1, "return self.#{name}")
        code << "}\n\n"
        code << "func (self *default#{cfg['name']}) Set#{camel}(#{name} time.Time) {\n"
        code << tabs(1, "if self.#{name} != #{name} {")
        code << tabs(2, "self.#{name} = #{name}")
        code << tabs(2, "self.updatedFields.Set(#{index + 1})")
        code << tabs(1, "}")
        code << "}\n\n"
        code << "func (self *default#{cfg['name']}) Set#{camel}Number(#{name} int) {\n"
        code << tabs(1, "self.#{name} = bsonmodel.NumberToDate(#{name})")
        code << tabs(1, "self.updatedFields.Set(#{index + 1})")
        code << "}\n\n"
      end
    when 'object'
      code << "func (self *default#{cfg['name']}) #{camel}() #{field['model']} {\n"
      code << tabs(1, "return self.#{name}")
      code << "}\n\n"
    when 'map'
      key_type = field['key']
      value_type = field['value']
      code << "func (self *default#{cfg['name']}) #{camel}() #{map_type(key_type)} {\n"
      code << tabs(1, "return self.#{name}")
      code << "}\n\n"
      if field.has_key? 'quick-access-method'
        code << "func (self *default#{cfg['name']}) #{field['quick-access-method']}(id #{key_type}) #{value_type} {\n"
        code << tabs(1, "value := self.#{name}.Get(id)")
        code << tabs(1, "if value == nil {")
        code << tabs(2, "return nil")
        code << tabs(1, "}")
        code << tabs(1, "return value.(#{value_type})")
        code << "}\n\n"
      elsif camel.end_with? 's'
        code << "func (self *default#{cfg['name']}) #{camel[0..-2]}(id #{key_type}) #{value_type} {\n"
        code << tabs(1, "value := self.#{name}.Get(id)")
        code << tabs(1, "if value == nil {")
        code << tabs(2, "return nil")
        code << tabs(1, "}")
        code << tabs(1, "return value.(#{value_type})")
        code << "}\n\n"
      end
    when 'simple-map'
      key_type = field['key']
      code << "func (self *default#{cfg['name']}) #{camel}() #{simple_map_type(key_type)} {\n"
      code << tabs(1, "return self.#{name}")
      code << "}\n\n"
    else
      raise "unsupported field type `#{field['type']}` on #{cfg['name']}.#{field['name']}"
    end
  end
end

def fill_new(code, cfg, has_parent = false)
  if has_parent
    code << "func New#{cfg['name']}(parent #{cfg['parent']['name']}) #{cfg['name']} {\n"
    code << tabs(1, "self := &default#{cfg['name']}{updatedFields: &bitset.BitSet{}, parent: parent}")
  else
    code << "func New#{cfg['name']}() #{cfg['name']} {\n"
    code << tabs(1, "self := &default#{cfg['name']}{updatedFields: &bitset.BitSet{}}")
  end
  cfg['fields'].each do |field|
    name = field['name']
    bname = field['bname']
    case field['type']
    when 'object'
      code << tabs(1, "self.#{name} = New#{field['model']}(self)")
    when 'map'
      key_type = field['key']
      value_type = field['value']
      code << tabs(1, "self.#{name} = #{map_factory(key_type)}(self, \"#{bname}\", #{value_type}Factory())")
    when 'simple-map'
      key_type = field['key']
      value_type = field['value']
      code << tabs(1, "self.#{name} = #{simple_map_factory(key_type)}(self, \"#{bname}\", #{simple_value_type(value_type)})")
    end
  end
  code << tabs(1, "return self")
  code << "}\n\n"
end

def fill_encoder(code, cfg)
  small_camel = to_small_camel(cfg['name'])
  # struct
  code << "type #{small_camel}Encoder struct{}\n\n"
  # IsEmpty
  code << "func (codec *#{small_camel}Encoder) IsEmpty(ptr unsafe.Pointer) bool {\n"
  code << tabs(1, "return false")
  code << "}\n\n"
  # Encode
  code << "func (codec *#{small_camel}Encoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {\n"
  code << tabs(1, "p := ((*default#{cfg['name']})(ptr))")
  code << tabs(1, "stream.WriteObjectStart()")
  first = true
  cfg['fields'].each do |field|
    next if field['json-ignore'] == true
    name = field['name']
    if first
      first = false
    else
      code << tabs(1, "stream.WriteMore()")
    end
    code << tabs(1, "stream.WriteObjectField(\"#{name}\")")
    case field['type']
    when 'int'
      if field['virtual'] == true
        code << tabs(1, "stream.WriteInt(p.#{to_camel(name)}())")
      else
        code << tabs(1, "stream.WriteInt(p.#{name})")
      end
    when 'string'
      if field['virtual'] == true
        code << tabs(1, "stream.WriteString(p.#{to_camel(name)}())")
      else
        code << tabs(1, "stream.WriteString(p.#{name})")
      end
    when 'float64'
      if field['virtual'] == true
        code << tabs(1, "stream.WriteFloat64(p.#{to_camel(name)}())")
      else
        code << tabs(1, "stream.WriteFloat64(p.#{name})")
      end
    when 'datetime'
      if field['virtual'] == true
        code << tabs(1, "stream.WriteInt64(p.#{to_camel(name)}().Unix())")
      else
        code << tabs(1, "stream.WriteInt64(p.#{name}.Unix())")
      end
    when 'date'
      if field['virtual'] == true
        code << tabs(1, "stream.WriteInt(bsonmodel.DateToNumber(p.#{to_camel(name)}()))")
      else
        code << tabs(1, "stream.WriteInt(bsonmodel.DateToNumber(p.#{name}))")
      end
    else
      code << tabs(1, "stream.WriteVal(p.#{name})")
    end
  end
  code << tabs(1, "stream.WriteObjectEnd()")
  code << "}\n\n"
  # init
  code << "func init() {\n"
  code << tabs(1, "jsoniter.RegisterTypeEncoder(\"#{cfg['package']}.default#{cfg['name']}\", &#{small_camel}Encoder{})")
  code << "}\n"
end


def generate_root(cfg)
  code = "package #{cfg['package']}\n\n"
  fill_imports(code, cfg)
  fill_interface(code, 'bsonmodel.RootModel', cfg)
  fill_const(code, cfg)
  fill_struct(code, cfg)
  fill_to_bson(code, cfg)
  fill_to_data(code, cfg)
  fill_load_jsoniter(code, cfg, true)
  fill_reset(code, cfg)
  fill_any_updated(code, cfg)
  fill_any_deleted(code, cfg)
  code << "func (self *default#{cfg['name']}) Parent() bsonmodel.BsonModel {\n"
  code << tabs(1, "return nil")
  code << "}\n\n"
  code << "func (self *default#{cfg['name']}) XPath() bsonmodel.DotNotation {\n"
  code << tabs(1, "return bsonmodel.RootPath()")
  code << "}\n\n"
  fill_append_updates(code, cfg)
  fill_to_document(code, cfg)
  fill_load_document(code, cfg, true)
  fill_deleted_size(code, cfg)
  fill_fully_update(code, cfg, true)
  fill_to_sync(code, cfg, true)
  fill_to_delete(code, cfg)
  code << "func (self *default#{cfg['name']}) ToUpdate() bson.M {\n"
  code << tabs(1, "if self.AnyUpdated() {")
  code << tabs(2, "return self.AppendUpdates(bson.M{})")
  code << tabs(1, "}")
  code << tabs(1, "return bson.M{}")
  code << "}\n\n"
  code << "func (self *default#{cfg['name']}) MarshalToJsonString() (string, error) {\n"
  code << tabs(1, "return jsoniter.MarshalToString(self)")
  code << "}\n\n"
  fill_xetters(code, cfg)
  fill_new(code, cfg)
  small_camel = to_small_camel(cfg['name'])
  code << "func Load#{cfg['name']}FromDocument(m bson.M) (#{small_camel} #{cfg['name']}, err error) {\n"
  code << tabs(1, "#{small_camel} = New#{cfg['name']}()")
  code << tabs(1, "err = #{small_camel}.LoadDocument(m)")
  code << tabs(1, "return")
  code << "}\n\n"
  code << "func Load#{cfg['name']}FromJsoniter(any jsoniter.Any) (#{small_camel} #{cfg['name']}, err error) {\n"
  code << tabs(1, "#{small_camel} = New#{cfg['name']}()")
  code << tabs(1, "err = #{small_camel}.LoadJsoniter(any)")
  code << tabs(1, "return")
  code << "}\n\n"
  fill_encoder(code, cfg)
  code << "\n"
end

def generate_object(cfg)
  code = "package #{cfg['package']}\n\n"
  fill_imports(code, cfg)
  fill_interface(code, 'bsonmodel.ObjectModel', cfg)
  fill_const(code, cfg)
  static_xpath = all_object(cfg)
  if static_xpath
    parent = cfg['parent']
    xpath = ["\"#{cfg['bname']}\""]
    p = parent
    until p['parent'].nil?
      xpath.insert(0, "\"#{p['bname']}\"")
      p = p['parent']
    end
    code << "var xpath#{cfg['name']} bsonmodel.DotNotation = bsonmodel.PathOfNames(#{xpath.join(', ')})\n\n"
  end
  fill_struct(code, cfg)
  fill_to_bson(code, cfg)
  fill_to_data(code, cfg)
  fill_load_jsoniter(code, cfg)
  fill_reset(code, cfg)
  fill_any_updated(code, cfg)
  fill_any_deleted(code, cfg)
  code << "func (self *default#{cfg['name']}) Parent() bsonmodel.BsonModel {\n"
  code << tabs(1, "return self.parent")
  code << "}\n\n"
  if static_xpath
    code << "func (self *default#{cfg['name']}) XPath() bsonmodel.DotNotation {\n"
    code << tabs(1, "return xpath#{cfg['name']}")
    code << "}\n\n"
  else
    code << "func (self *default#{cfg['name']}) XPath() bsonmodel.DotNotation {\n"
    code << tabs(1, "return self.parent.XPath().Resolve(#{cfg['bname']})")
    code << "}\n\n"
  end
  fill_append_updates(code, cfg)
  fill_to_document(code, cfg)
  fill_load_document(code, cfg)
  fill_deleted_size(code, cfg)
  fill_fully_update(code, cfg)
  fill_to_sync(code, cfg)
  fill_to_delete(code, cfg)
  fill_xetters(code, cfg)
  fill_new(code, cfg, true)
  fill_encoder(code, cfg)
  code << "\n"
end

def generate_map_value(cfg)
  code = "package #{cfg['package']}\n\n"
  fill_imports(code, cfg)
  key_type = cfg['key']
  fill_interface(code, map_value_type(key_type), cfg)
  fill_const(code, cfg)
  fill_struct(code, cfg, map_value_struct(key_type))
  fill_to_bson(code, cfg)
  fill_to_data(code, cfg)
  fill_load_jsoniter(code, cfg)
  fill_reset(code, cfg)
  fill_any_updated(code, cfg)
  fill_any_deleted(code, cfg)
  fill_append_updates(code, cfg)
  fill_to_document(code, cfg)
  fill_load_document(code, cfg)
  fill_deleted_size(code, cfg)
  fill_fully_update(code, cfg)
  fill_to_sync(code, cfg)
  fill_to_delete(code, cfg)
  fill_xetters(code, cfg)
  fill_new(code, cfg)
  small_camel = to_small_camel(cfg['name'])
  code << "var #{small_camel}Factory #{map_value_factory(key_type)} = func() #{map_value_type(key_type)} {\n"
  code << tabs(1, "return New#{cfg['name']}()")
  code << "}\n\n"
  code << "func #{cfg['name']}Factory() #{map_value_factory(key_type)} {\n"
  code << tabs(1, "return #{small_camel}Factory")
  code << "}\n\n"
  fill_encoder(code, cfg)
  code << "\n"
end


cfg = File.open(ARGV[0]) { |io| YAML.load io.read }

parents = Hash.new
bnames = Hash.new
map_models = Set.new

unless cfg.has_key? 'package'
  raise "missing required field `package`"
end

cfg['objects'].each do |model|
  model['package'] = cfg['package']
  unless model.has_key? 'name'
    raise "missing required field `name`"
  end
  unless model.has_key? 'file'
    model['file'] = "#{to_underscore(model['name'])}.go"
  end
  unless model['file'].end_with? '.go'
    model['file'] = "#{model['file']}.go"
  end
  model['fields'].each_with_index do |field, index|
    if field['type'] == 'object'
      parents[field['model']] = model
      bnames[field['model']] = field['bname']
    end
    unless field.has_key? 'bname'
      field['bname'] = field['name']
    end
    if field['virtual']
      model['fields'].select do |v| 
        field['sources'].include?(v['name'])
      end.each do |v|
        if v.has_key? 'relations'
          v['relations'] << index + 1
        else
          v['relations'] = [index + 1]
        end
      end
    end
    if field['type'] == 'map'
      map_models << field['value']
    end
  end
end.each do |model|
  name = model['name']
  if parents.has_key? name
    model['parent'] = parents[name]
    model['bname'] = bnames[name]
  end
  if map_models.include? name
    model['type'] = 'map-value'
  end
end.each do |model|
  case model['type']
  when 'root'
    code = generate_root(model)
  when 'object'
    code = generate_object(model)
  when 'map-value'
    code = generate_map_value(model)
  else
    raise "unknown type #{model['type']}"
  end
  require 'fileutils'
  package_dir = File.join(ARGV[1], model['package'])
  unless File.directory?(package_dir)
    FileUtils.mkdir_p(package_dir)
  end
  filename = model['file']
  puts "Generating #{filename} ... (on path: #{package_dir})"
  File.open(File.join(package_dir, filename), "w") do |io|
    io.syswrite(code)
  end
  puts "OK"
end

puts "Done."
