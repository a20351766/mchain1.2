Coding guidelines
-----------------

Coding Golang
~~~~~~~~~~~~~~

We code in Go™ and strictly follow the `best
practices <https://golang.org/doc/effective_go.html>`__ and will not
accept any deviations. You must run the following tools against your Go
code and fix all errors and warnings: -
`golint <https://github.com/golang/lint>`__ - `go
vet <https://golang.org/cmd/vet/>`__ -
`goimports <https://godoc.org/golang.org/x/tools/cmd/goimports>`__

API Documentation
^^^^^^^^^^^^^^^^^

The API documentation for Hyperledger Mchain's Golang APIs is available
in `GoDoc <https://godoc.org/github.com/hyperledger/mchain>`_.


Generating gRPC code
---------------------

If you modify any ``.proto`` files, run the following command to
generate/update the respective ``.pb.go`` files.

::

    cd $GOPATH/src/github.com/hyperledger/mchain
    make protos

Adding or updating Go packages
------------------------------

Hyperledger Mchain uses Go Vendoring for package
management. This means that all required packages reside in the
``$GOPATH/src/github.com/hyperledger/mchain/vendor`` folder. Go will use
packages in this folder instead of the GOPATH when the ``go install`` or
``go build`` commands are executed. To manage the packages in the
``vendor`` folder, we use
`dep <https://golang.github.io/dep/>`__.

.. Licensed under Creative Commons Attribution 4.0 International License
   https://creativecommons.org/licenses/by/4.0/
