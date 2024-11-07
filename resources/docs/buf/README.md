# Buf Build Configuration

Buf build configuration is stored in `buf.yaml` in the repository.

- buf.yaml

```yaml
version: v2
# The v2 buf.yaml file specifies a local workspace, which consists of at least one module.
# The buf.yaml file should be placed at the root directory of the workspace, which
# should generally be the root of your source control repository.
modules:
  # Each module entry defines a path, which must be relative to the directory where the
  # buf.yaml is located. You can also specify directories to exclude from a module.
  - path: proto/foo
    # Modules can also optionally specify their Buf Schema Repository name if it exists.
    name: buf.build/acme/foo
    # Excluding a subdirectory and a specific .proto file. Note that the paths for exclusion
    # are relative to the buf.yaml file.
  - path: proto/bar
    name: buf.build/acme/bar
    excludes:
      - proto/bar/a
      - proto/bar/b/e.proto
    # A module can have its own lint and breaking configuration, which overrides the default
    # lint and breaking configuration in its entirety for that module. All values from the
    # default configuration are overridden and no rules are merged.
    lint:
      use:
        - STANDARD
      except:
        - IMPORT_USED
      ignore:
        - proto/bar/c
      ignore_only:
        ENUM_ZERO_VALUE_SUFFIX:
          - proto/bar/a
          - proto/bar/b/f.proto
    # v1 configurations had an allow_comment_ignores option to opt-in to comment ignores.
    #
    # In v2, we allow comment ignores by default, and allow opt-out from comment ignores
    # with the disallow_comment_ignores option.
    disallow_comment_ignores: false
    enum_zero_value_suffix: _UNSPECIFIED
    rpc_allow_same_request_response: false
    rpc_allow_google_protobuf_empty_requests: false
    rpc_allow_google_protobuf_empty_responses: false
    service_suffix: Service
    disable_builtin: false
    # Breaking configuration for this module only. Behaves the same as a module-level
    # lint configuration.
    breaking:
      use:
        - FILE
      except:
        - EXTENSION_MESSAGE_NO_DELETE
      ignore_unstable_packages: true
      disable_builtin: false
  # Multiple modules are allowed to have the same path, as long as they don't share '.proto' files.
  - path: proto/common
    module: buf.build/acme/weather
    includes:
      - proto/common/weather
  - path: proto/common
    module: buf.build/acme/location
    includes:
      - proto/common/location
    excludes:
      # Excludes and includes can be specified at the same time, but if they are, each directory
      # in excludes must be contained in a directory in includes.
      - proto/common/location/test
  - path: proto/common
    module: buf.build/acme/other
    excludes:
      - proto/common/location
      - proto/common/weather
# Dependencies shared by all modules in the workspace. Must be modules hosted in the Buf Schema Registry.
# The resolution of these dependencies is stored in the buf.lock file.
deps:
  - buf.build/acme/paymentapis # The latest accepted commit.
  - buf.build/acme/pkg:47b927cbb41c4fdea1292bafadb8976f # The '47b927cbb41c4fdea1292bafadb8976f' commit.
  - buf.build/googleapis/googleapis:v1beta1.1.0 # The 'v1beta1.1.0' label.
# The default lint configuration for any modules that don't have a specific lint configuration.
#
# If this section isn't present, AND a module doesn't have a specific lint configuration, the default
# lint configuration is used for the module.
lint:
  use:
    - STANDARD
    - TIMESTAMP_SUFFIX # This rule comes from the plugin example below.
# Default breaking configuration. It behaves the same as the default lint configuration.
breaking:
  use:
    - FILE
# Optional Buf plugins. Can use to require custom lint or breaking change rules specified in a locally
# installed plugin. Each Buf plugin is listed separately, and can include options if the plugin allows
# for them. If a rule has its `default` value set to true, the rule will be checked against even if
# the 'lint' and 'breaking' fields aren't set.
#
# See the example at https://github.com/bufbuild/bufplugin-go/blob/main/check/internal/example/cmd/buf-plugin-timestamp-suffix/main.go
# for more detail about the sample below.
plugins:
  - plugin: plugin-timestamp-suffix # Specifies the installed plugin to use
    options:
      # The TIMESTAMP_SUFFIX rule specified above allows the user to change the suffix by providing a
      # new value here.
      timestamp_suffix: _time
```

- buf.gen.yaml

```yaml
version: v2
# 'clean', when set to true, deletes the directories, zip files, and/or jar files specified in the `out` field for
# all plugins before running code generation.
clean: true
# 'managed' contains the configuration for managed mode: https://buf.build/docs/generate/managed-mode
# It has three top-level keys: 'enabled', 'disable', and 'override'.
#
# When managed mode is enabled, it uses default values for certain file and field options during code
# generation. Options, accepted values, and defaults are documented here:
# https://buf.build/docs/generate/managed-mode#default-behavior
# The 'disable' key configures modules, paths, fields, and/or options that are excluded from managed
# mode's behavior. The 'override' key configures field and file option values that override the
# default values managed mode uses during code generation.
#
# In the case of options that combine with other options (e.g. java_package + java_package_prefix
# + java_package_suffix), they are all applied if possible. If not (e.g. when all three are set)
# then the last configuration rule wins.
managed:
  # 'enabled: true' turns managed mode on, 'enabled: false' ignores all managed mode options.
  enabled: true # default: false
  # 'disable' is a list of 'disable' rules managing either file options or field options.
  # A 'disable' rule must have at least one key set.
  disable:
    # Don't modify any files in buf.build/googleapis/googleapis
    - module: buf.build/googleapis/googleapis
    # Don't modify any files in the foo/v1 directory. This can be a path to a directory
    # or a .proto file. If it's a directory path, all .proto files in the directory are
    # ignored.
    - path: foo/v1
    # Ignore the csharp_namespace file option for all modules and files in the input.
    - file_option: csharp_namespace
    # Ignore the js_type field option for any file.
    - field_option: js_type
    # Ignore the foo.bar.Baz.field_name field.
    - field: foo.bar.Baz.field_name
    # Setting all 3 for file options: don't modify java_package and go_package (by disabling
    # setting the go_package_prefix) for files in foo/v1 in buf.build/acme/weather
    - module: buf.build/acme/weather
      path: foo/v1
      file_option: java_package
    # Setting all 4 for field options: disable js_type for all files that match the
    # module, path, and field name.
    - module: buf.build/acme/petapis
      field: foo.bar.Baz.field_name
      path: foo/v1
      field_option: js_type
  # 'override' is a list of 'override' rules for the list of field and file options that
  # managed mode handles.
  override:
    # When 'file_option' and 'value' are set, managed mode uses the value set with
    # this rule instead of the default value or sets it to this value if the default
    # value is none.
    #
    # Example: Modify the java_package options to <net>.<proto_package> for all files.
    # Also modify go_package_prefix to be company.com/foo/bar for all files.
    - file_option: java_package_prefix
      value: net
    - file_option: go_package_prefix
      value: company.com/foo/bar
    # When 'file_option', 'value', and 'module' are set, managed mode uses the value
    # set in this rule instead of the default value for all files in the specified module.
    # It sets it to this value if the default value is none.
    #
    # Example: Modify the java_package option to <com>.<proto_package>.<com> for all files
    # in buf.build/acme/petapis. Also modify go_package_prefix to be company.com/foo/baz
    # for all files in buf.build/acme/petapis.
    # These rules take precedence over the rule above.
    - file_option: java_package_prefix
      module: buf.build/acme/petapis
      value: com
    - file_option: java_package_suffix
      module: buf.build/acme/petapis
      value: com
    - file_option: go_package_prefix
      module: buf.build/acme/petapis
      value: company.com/foo/baz
    # When 'file_option', 'value', and 'path' are set, managed mode uses the value set
    # in this rule instead of the default value for the specific file path. If the path
    # is a directory, the rule affects all .proto files in the directory. Otherwise, it
    # only affects the specified .proto file.
    #
    # Example: For the file foo/bar/baz.proto, set java_package specifically to
    # "com.x.y.z". Also for the file foo/bar/baz.proto modify go_package_prefix to be
    # company.com/bar/baz.
    # This takes precedence over the previous rules above.
    - file_option: java_package
      path: foo/bar/baz.proto
      value: com.x.y.z
    - file_option: go_package_prefix
      path: foo/bar/baz.proto
      value: company.com/bar/baz
    # When 'field_option', 'value', and 'module' are set, managed mode uses the value
    # set in this rule instead of the default value for all files in the specified module.
    # It sets it to this value if the default value is none.
    #
    # Example: For all fields in the buf.build/acme/paymentapis module where the field is
    # one of the compatible types, set the 'js_type' to "JS_NORMAL".
    - field_option: js_type
      module: buf.build/acme/paymentapis
      value: JS_NORMAL
    # When 'field_option', 'value', and 'field' are set, managed mode uses the value set
    # in this rule instead of the default value of the specified field. It sets it to
    # this value if the default value is none.
    #
    # Example: Set the package1.Message2.field3 field 'js_type' value to "JS_STRING". You can
    # additionally specify the module and path, but the field name is sufficient.
    - field_option: js_type
      value: JS_STRING
      field: package1.Message2.field3
# 'plugins' is a list of plugin configurations used for buf generate.
#
# A 'plugin' configuration has 8 possible keys:
#  - one of (required):
#    - 'remote': remote plugin name (e.g. buf.build/protocolbuffers/go)
#    - 'protoc_builtin': a 'protoc' built-in plugin (e.g. 'cpp' for 'protoc-gen-cpp')
#    - 'local': a string or list of strings that point to a protoc plugin binary on your
#      '${PATH}'. If a list of strings is specified, the first is the binary name, and the
#      subsequent strings are considered arguments passed to the binary.
#  - 'out': <string> path to the file output, which is the same as v1 (required)
#  - 'opt': a list of plugin options, which is the same as v1 (optional)
#  - 'strategy': a string for the invocation strategy, which is the same as v1 (optional)
#  - 'include_imports': <boolean> (optional, precedence given to CLI flag)
#  - 'include_wkt': <boolean> (optional, precedence given to CLI flag)
plugins:
  # BSR remote plugin
  - remote: buf.build/protocolbuffers/go
    out: gen/proto
  # Built-in protoc plugin for C++
  - protoc_builtin: cpp
    protoc_path: /path/to/protoc
    out: gen/proto
  # Local binary plugin, search in ${PATH} by default
  - local: protoc-gen-validate
    out: gen/proto
  # Relative paths automatically work
  - local: path/to/protoc-gen-validate
    out: gen/proto
  # Absolute paths automatically work
  - local: /usr/bin/path/to/protoc-gen-validate
    out: gen/proto
  # Binary plugin with arguments and includes
  - local: [ "go", "run", "google.golang.org/protobuf/cmd/protoc-gen-go" ]
    out: gen/proto
    opt:
      - paths=source_relative
      - foo=bar
      - baz
    strategy: all
    include_imports: true
    include_wkt: true
# 'inputs' is a list of inputs that will be run for buf generate. It's a
# required key for v2 buf.gen.yaml and allows you to specify options based on the type
# of input (https://buf.build/docs/reference/inputs.md#source) being configured.
inputs:
  # Git repository
  - git_repo: github.com/acme/weather
    branch: dev
    subdir: proto
    depth: 30
  # BSR module with types and include, and exclude keys specified
  - module: buf.build/acme/weather:main
    types:
      - "foo.v1.User"
      - "foo.v1.UserService"
    # If empty, include all paths.
    paths:
      - a/b/c
      - a/b/d
    exclude_paths:
      - a/b/c/x.proto
      - a/b/d/y.proto
  # Local module at provided directory path
  - directory: x/y/z
  # Tarball at provided directory path. Automatically derives compression algorithm from
  # the file extension:
  # - '.tgz' and '.tar.gz' extensions automatically use Gzip
  # - '.tar.zst' automatically uses Zstandard.
  - tarball: a/b/x.tar.gz
  # Tarball with 'compression', 'strip_components', and 'subdir' keys set explicitly.
  # - 'strip_components=<int>' reads at the relative path and strips some number of
  #    components—in this example, 2.
  # - 'subdir=<string>' reads at the relative path and uses the subdirectory specified
  #    within the archive as the base directory—in this case, 'proto'.
  - tarball: c/d/x.tar.zst
    compression: zstd
    strip_components: 2
    subdir: proto
  # The same applies to 'zip' inputs. A 'zip' input is read at the relative path or http
  # location, and can set 'strip_components' and 'subdir' optionally.
  - zip_archive: https://github.com/googleapis/googleapis/archive/master.zip
    strip_components: 1
  # 'proto_file' is a path to a specific proto file. Optionally, you can include package
  # files as part of the files to be generated (the default is false).
  - proto_file: foo/bar/baz.proto
    include_package_files: true
  # We also support Buf images as inputs. Images can be any of the following formats:
  #  - 'binary_image'
  #  - 'json_image'
  #  - 'txt_image'
  #  - 'yaml_image'
  # Each image format also supports compression optionally.
  #
  # The example below is a binary Buf image with compression set for Gzip.
  - binary_image: image.binpb.gz
    compression: gzip
```

## Update dependencies

```shell
$ buf dep update
```

- buf.lock

```yaml
# Generated by buf. DO NOT EDIT.
version: v2
deps:
  - name: buf.build/googleapis/googleapis
    commit: 7a6bc1e3207144b38e9066861e1de0ff
    digest: b5:6d05bde5ed4cd22531d7ca6467feb828d2dc45cc9de12ce3345fbddd64ddb1bf0db756558c32ca49e6bc7de4426ada8960d5590e8446854b81f5f36f0916dc48
```

See [buf.build](https://buf.build) for more details