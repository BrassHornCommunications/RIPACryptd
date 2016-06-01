// Copyright (c) 2016 Brass Horn Communications
// Use of this sourcecode is governed by the BSD
// 2 Clause License (see LICENSE.md)

/*
RIPACryptd is the server side daemon of https://RIPACrypt.download - A SaaS
platform that stores secrets unless a series of checkins are missed at which
point the secrets are destroyed.

This is a thought experiment in defeating the UK's compelled decryption laws.

We strongly recommend using the https://RIPACrypt.download server rather than
running your own instance but the decision is yours.

The client can be found at github.com/BrassHornCommunications/RIPACrypt

Usage:
ripacryptd [OPTIONS]

Application Options:
	-conf				Specify the path to a configuration file

*/

package main
