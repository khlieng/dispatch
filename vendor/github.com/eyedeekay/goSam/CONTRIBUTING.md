How to make contributions to goSam
==================================

Welcome to goSam, the easy-to-use http client for i2p. We're glad you're here
and interested in contributing. Here's some help getting started.

Table of Contents
-----------------

  * (1) Environment
  * (2) Testing
  * (3) Filing Issues/Reporting Bugs/Making Suggestions
  * (4) Contributing Code/Style Guide
    - (a) Adding i2cp and tunnel Options
    - (b) Writing Tests
    - (c) Style
    - (d) Other kinds of modification?
  * (5) Conduct

### (1) Environment

goSam is a simple go library. You are free to use an IDE if you wish, but all
that is required to build and test the library are a go compiler and the gofmt
tool. Git is the version control system. All the files in the library are in a
single root directory. Invoking go build from this directory not generate any
files.

### (2) Testing

Tests are implemented using the standard go "testing" library in files named
"file\_test.go," so tests of the client go in client\_test.go, name lookups
in naming\_test.go, et cetera. Everything that can be tested, should be tested.

Testing is done by running

        go test

More information about designing tests is below in the
**Contributing Code/Style Guide** section below.

### (3) Filing issues/Reporting bugs/Making suggestions

If you discover the library doing something you don't think is right, please let
us know! Just filing an issue here is OK.

If you need to suggest a feature, we're happy to hear from you too. Filing an
issue will give us a place to discuss how it's implemented openly and publicly.

Please file an issue for your new code contributions in order to provide us with
a place to discuss them for inclusion.

### (4) Contributing Code/Style Guide

Welcome new coders. We have good news for you, this library is really easy to
contribute to. The easiest contributions take the form of i2cp and tunnel
options.

#### (a) Adding i2cp and tunnel Options

First, add a variable to store the state of your new option. For example, the
existing variables are in the Client class [here:](https://github.com/cryptix/goSam/blob/701d7fcf03ddb354262fe213163dcf6f202a24f1/client.go#L29)

i2cp and tunnel options are added in a highly uniform process of basically three
steps. First, you create a functional argument in the options.go file, in the
form:

``` Go
        // SetOPTION sets $OPTION
        func SetOPTION(arg type) func(*Client) error {  // arg type
            return func(c *Client) error {              // pass a client to the inner function and declare error return function
                if arg == valid {                       // validate the argument
                    c.option = s                        // set the variable to the argument value
                    return nil                          // if option is set successfully return nil error
                }
                return fmt.Errorf("Invalid argument:" arg) // return a descriptive error if arg is invalid
            }
        }
```

[example](https://github.com/cryptix/goSam/blob/701d7fcf03ddb354262fe213163dcf6f202a24f1/options.go#L187)

Next, you create a getter which prepares the option. Regardless of the type of
option that is set, these must return strings representing valid i2cp options.

``` Go
        //return the OPTION as a string.
        func (c *Client) option() string {
            return fmt.Sprintf("i2cp.option=%d", c.option)
        }
```

[example](https://github.com/cryptix/goSam/blob/701d7fcf03ddb354262fe213163dcf6f202a24f1/options.go#L299)

Lastly, you'll need to add it to the allOptions function and the
Client.NewClient() function. To add it to allOptions, it looks like this:

``` Go
        //return all options as string ready for passing to sendcmd
        func (c *Client) allOptions() string {
            return c.inlength() + " " +
                c.outlength() + " " +
                ... //other options removed from example for brevity
                c.option()
        }
```

``` Go
        //return all options as string ready for passing to sendcmd
        func (c *Client) NewClient() (*Client, error) {
            return NewClientFromOptions(
                SetHost(c.host),
                SetPort(c.port),
                ... //other options removed from example for brevity
                SetCompression(c.compression),
                setlastaddr(c.lastaddr),
                setid(c.id),
            )
        }
```

[example](https://github.com/cryptix/goSam/blob/701d7fcf03ddb354262fe213163dcf6f202a24f1/options.go#L333)

#### (b) Writing Tests

Before the feature can be added, you'll need to add a test for it to
options_test.go. To do this, just add your new option to the long TestOptions
functions in options_test.go.

``` Go
        func TestOptionHost(t *testing.T) {
            client, err := NewClientFromOptions(
                SetHost("127.0.0.1"),
                SetPort("7656"),
                ... //other options removed from example for brevity
                SetCloseIdleTime(300001),
            )
            if err != nil {
                t.Fatalf("NewClientFromOptions() Error: %q\n", err)
            }
            if result, err := client.validCreate(); err != nil {
                t.Fatalf(err.Error())
            } else {
                t.Log(result)
            }
            client.CreateStreamSession("")
            if err := client.Close(); err != nil {
                t.Fatalf("client.Close() Error: %q\n", err)
            }
        }

        func TestOptionPortInt(t *testing.T) {
            client, err := NewClientFromOptions(
                SetHost("127.0.0.1"),
                SetPortInt(7656),
                ... //other options removed from example for brevity
                SetUnpublished(true),
            )
            if err != nil {
                t.Fatalf("NewClientFromOptions() Error: %q\n", err)
            }
            if result, err := client.validCreate(); err != nil {
                t.Fatalf(err.Error())
            } else {
                t.Log(result)
            }
            client.CreateStreamSession("")
            if err := client.Close(); err != nil {
                t.Fatalf("client.Close() Error: %q\n", err)
            }
        }

```

If any of these tasks fail, then the test should fail.

#### (c) Style

It's pretty simple to make sure the code style is right, just run gofmt over it
to adjust the indentation, and golint over it to ensure that your comments are
of the correct form for the documentation generator.

#### (d) Other kinds of modification?

It may be useful to extend goSam in other ways. Since there's not a
one-size-fits-all uniform way of dealing with these kinds of changes, open an
issue for discussion and

### (5) Conduct

This is a small-ish, straightforward library intended to enable a clear
technical task. We should be able to be civil with eachother, and give and
accept criticism contructively and respectfully.

This document was drawn from the examples given by Mozilla
[here](mozillascience.github.io/working-open-workshop/contributing/)
