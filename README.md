<h1 align="center">Radon</h1>

This is a total rewrite of Radon, with more focus on code quality, testing, and some slight improvements to the actual language.

For syntax highlighting and snippets in [TextMate](http://macromates.com), install [Radon.tmbundle](Radon.tmbundle) from this repository.

Here's a nice example of some of the features:

```
calc input = do
	result = 0

	for character in input do
		result = match character where
			| "+" -> result + 1,
			| "-" -> result - 1,
			| "*" -> result * 2,
			| "/" -> result / 2
	end
end
```
