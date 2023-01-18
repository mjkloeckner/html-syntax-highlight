# Highlight code on a static html file

This program parses an html file and adds
[chroma](https://github.com/alecthomas/chroma) classes to code within pre and
code tags. Then you can add chroma themes with css.

## Running

Make sure you have [golang](https://go.dev/) installed (if not follow this
[instructions](https://go.dev/doc/install))then run:

```shell
$ go run main.go <test_file.html>
```

with `<test_file.html>` being an html file containing some code to highlight.
This is an example of a code snippet within an html file that this program
should add css classes

```html
(...)
<pre><code class="language-c">
#include <stdio.h>

int main(void) {
    char *hw = "Hello, world!\n";

    printf("%s\n", hw);
    return 0;
}
</code></pre>
(...)
```

But it uses regular expression so it should be easier to adapt to other formats.

## Example
The main purpose of this program was to highlight code snippets on my blog.
Checkout [this post](https://kloeckner.com.ar/blog/testing-syntax-highlight/testing-syntax-highlight)
that I made to test it.

## License
[MIT](https://opensource.org/licenses/MIT)
