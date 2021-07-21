# Metadata for Paketo Build/Run Stack Images

The following build/run stack images are supported:\
**Base**:\
Build Images
- `docker.io/paketobuildpacks/build:base-cnb`

Run Images
- `docker.io/paketobuildpacks/run:base-cnb`
- `gcr.io/paketo-buildpacks/run:base-cnb`

**Full**:\
Build Images
- `docker.io/paketobuildpacks/build:full-cnb`

Run Images
- `docker.io/paketobuildpacks/run:full-cnb`
- `gcr.io/paketo-buildpacks/run:full-cnb`

**Tiny**:\
Build Images
- `docker.io/paketobuildpacks/build:tiny-cnb`

Run Images
- `docker.io/paketobuildpacks/run:tiny-cnb`
- `gcr.io/paketo-buildpacks/run:tiny-cnb`

**Note:** These images are tagged in the format `<stack-name>-cnb`. For many of them, there is a corresponding image tagged `<stack-name>` (i.e. `paketobuildpacks/build:base`). Those images are useful for extending our stack with your own packages and metadata. The images with the `<stack-name>-cnb` already have [CNB metadata](https://github.com/buildpacks/spec/blob/main/platform.md#stacks) set and can be used directly as stack images.

## Use Cases


### Base (aka "bionic")
Ideal for:
- .NET Core apps
- Java apps and Go apps that require some C libraries
- Node.js/Python/Ruby/etc. apps **without** many native extensions

Contains:
- Build: ubuntu:bionic + openssl + CA certs + compilers + shell utilities
- Run: ubuntu:bionic + openssl + CA certs

### Full
Ideal for:
- PHP/Node.js/Python/Ruby/etc. apps **with** many native extensions

Contains:
- Build: ubuntu:bionic + many common C libraries and utilities
- Run: ubuntu:bionic + many common libraries and utilities

### Tiny
Ideal for:
- Most Go apps
- Java apps and Java [GraalVM Native Images](https://www.graalvm.org/docs/reference-manual/native-image/)

Contains:
- Build: ubuntu:bionic + openssl + CA certs + compilers + shell utilities
- Run: distroless-like bionic + glibc + openssl + CA certs
