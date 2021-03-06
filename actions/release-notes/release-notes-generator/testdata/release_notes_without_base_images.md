## Release Images

### Runtime Base Images

#### For CNB Builds:
- Tag: **`some-registry/run:1.0-tiny-cnb`**
- Digest: `some-registry/run@sha256:some-cnb-sha`

### Build-time Base Images

#### For CNB Builds:
- Tag: **`some-registry/build:1.0-tiny-cnb`**
- Digest: `some-registry/build@sha256:some-cnb-sha`
## Patched USNs
[USN-4498-1](https://ubuntu.com/security/notices/USN-4498-1):  Loofah vulnerability
* [CVE-2019-15587](https://people.canonical.com/~ubuntu-security/cve/CVE-2019-15587): In the Loofah gem for Ruby through v2.3.0 unsanitized JavaScript may occur in sanitized output when a crafted SVG element is republished.

[USN-4593-1](https://ubuntu.com/security/notices/USN-4593-1):  FreeType vulnerability
* [CVE-2020-15999](https://people.canonical.com/~ubuntu-security/cve/CVE-2020-15999): A buffer overflow was discovered in Load_SBit_Png.

## Build Receipt Diff
```
-ii  ruby-loofah         1.6.10ubuntu0.1   amd64  some description
+ii  ruby-loofah         1.6.12ubuntu0.1   amd64  some description
-ii  libfreetype6:amd64  2.8.1-2ubuntu2    amd64  some other description
+ii  libfreetype6:amd64  2.8.1-2ubuntu2.1  amd64  some other description
```
