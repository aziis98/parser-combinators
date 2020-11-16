 # Notes

```

1 + 2 + 3 	-> (1 + (2 + 3))

1 * 2 * 3 + 4 	-> ((1 * (2 * 3)) +  3)

1 * 2 ^ 3 + 4 	-> ((1 * (2 ^ 3)) +  3)

1 * 2 ^ 3 ^ 2 + 4 	-> ((1 * (2 ^ (3 ^ 2))) +  3)

1 + 2 - 3 	-> (1 + (2 - 3))

2 * 3 / 4 * 5 	-> (2 * (3 / (4 * 5)))
				-> (((2 * 3) / 4) * 5)

1 + 2 + 3 + 4 	-> (((1 + 2) + 3) + 4)

----------------------------

// "1 + 2 + 3 + 4" -> @Plus(@Plus(@Plus(1, 2), 3), 4)

parseLeftAssoc() {
	parseSuffix() {
		expect("+")
		parseNumber() -> n2

		@PartialPlus(n2) ->
	}

	parseNumber() -> n1
	
	while peek() = "+"  {
		parseSuffix() -> @PartialPlus(n2)
		@Plus(plus, n2) -> plus
	}

	@Plus(n1, n2) ->
}

```


```grammar
S <- A | B
A <- x A y | x z y 
B <- x B y y | x z y y 
```



















