# json-mp

JSON &lt;-> MessagePack converter

#### To clone the repository, use the following command:

```sh
git clone https://github.com/romanzac/json-mp.git
```

#### Build prerequisites

- Go v1.20
- GNU Make

#### Build project

```sh
cd json-mp
make
```

#### Run unit tests for ecoding/decoding core (optional)

```sh
make test
```

#### Prepare shape file:

Program requires to define shape of the data structure which will be encoded in MessagePack format.
Please edit shape/shape.go with a type named "DataShape".

#### Run json-mp:

Encode JSON -> MessagePack

```sh
e.g.: ./json-mp -i data/sample.json -o data/sample.mp
```

Decode MessagePack -> JSON

```sh
e.g.: ./json-mp -d -i data/sample.mp -o data/sample_out.json
```

#### Supported JSON data types:

- Null, Bool, Number, String, Array, Object

#### Supported Go data types:

- Nil, Bool, Int, uInt, Float, String, Map, Slice, Array, Struct, Interface, Pointer

#### Usage guide:
```sh
./json-mp -h
Encodes file from JSON to MessagePack format (default mode)

Usage:
  json-mp [flags]

Flags:
  -d, --decode          decodes MessagePack to JSON format
  -h, --help            help for json-mp
  -i, --input string    input file path
  -o, --output string   output file path
  -v, --version         version for json-mp
```

#### Further development ideas:

- Streaming implementation
- Automatic shape generation and caching like at https://transform.tools/json-to-go
