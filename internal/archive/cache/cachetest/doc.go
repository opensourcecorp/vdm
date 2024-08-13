// Package cachetest exists because Go won't let me have a test for
// [cache.AddRemote] in the cache package because it imports [remotes.Git],
// which causes an import cycle. So, the tests for that are here instead.
package cachetest
