Classification Trees in Go
==========================

The Go Package ``github.com/seehuhn/classification`` implements
classification trees as described in [HaTiFrie09]_

Copyright (C) 2014, 2015  Jochen Voss

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

.. [HaTiFrie09] Trevor Hastie, Robert Tibshirani, Jerome Friedman:
	       *The Elements of Statistical Learning*, second
	       edition,Springer, 2009

Current Status
--------------

The package is still experimental and changes frequently.
(Constructive) comments on the API or the implementation would be most
welcome.

Installation
------------

This package can be installed using the ``go get`` command::

    go get github.com/seehuhn/classification

Usage
-----

The current interface can be seen via the package's onlinehelp, either
on godoc.org_ or on the command line::

    godoc github.com/seehuhn/classification

.. _godoc.org: http://godoc.org/github.com/seehuhn/classification

Online Resources
----------------

There are many online resources about classification trees available.
Some examples:

  * http://web.stanford.edu/~hastie/local.ftp/Springer/OLD/ESLII_print4.pdf
  * http://www.ise.bgu.ac.il/faculty/liorr/hbchap9.pdf
