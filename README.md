# Chanty

## Overview
Chanty is a proof-of-concept micro-framework for building actors in Go.

## Why Actors?
Go is a ~~great~~ language for building concurrent systems. However, it can be difficult to reason about the interactions between goroutines.\
The actor model provides a way to encapsulate state and behavior in a single unit, which can make it easier to reason about concurrent systems.

## Why?
I wanted to explore generics in Go, and I thought that building a framework which promotes lexical confinement of goroutines would be a fun way to do that.
...and also to flex, if we're being honest 😎
