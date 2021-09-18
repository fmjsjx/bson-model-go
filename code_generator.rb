require 'set'
require 'yaml'


cfg = File.open(ARGV[0]) { |io| YAML.load io.read }

parents = Hash.new
bnames = Hash.new
