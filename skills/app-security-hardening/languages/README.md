# Languages

Use this directory only for language-specific or framework-runtime-specific pitfalls.

If a lesson is really about API design, token lifecycle, logging, storage, or release controls, put the principle in `../controls/` and link here only for language/runtime details. Concrete how-to steps belong in `../implementation/`.

| File | Scope |
| --- | --- |
| `dart.md` | Dart and Flutter-specific concerns. |
| `kotlin-java.md` | Kotlin/Java Android-specific code patterns. |
| `swift.md` | Swift/iOS-specific code patterns. |
| `typescript.md` | TypeScript frontend/backend client code concerns. |

When language-specific guidance changes how engineers implement a control, update or verify the matching file in [`../implementation/`](../implementation/) in the same change.
