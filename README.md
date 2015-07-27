goliner
=======

Execute some Go one-liners.

Examples
--------

Trivial helloworld:

`$ goliner 'println("hello world!")'`

`hello world!`

More sophisticated example with modules usage:

`$ goliner 'fmt.Println(strings.Join([]string{"foo", "bar"}, ", "))'`

`foo, bar`
