# bittorrent

The perennial favorite programming pastime of yore.

## Usage

```
$ go build

$ ./bittorrent [torrent_file.torrent]
```

If no torrent file is provided then the default test file will be used, which is for the latest Debian ISO.

## Background

Implementation follows the [original BitTorrent spec][og_spec] and the [unofficial spec][other_spec] which has more details relevant for implementation.

## Notes

- Do requests to the tracker need to be truthful? Can we lie about how much we have downloaded/uploaded?
    - No there doesn't seem to be a way to enforce honesty. More details about exploits and dishonesty in the protocol are covered [here][blackhat]

## TODO

- Would be interesting to implement an algorithm from the economics of BitTorrent paper: http://bittorrent.org/bittorrentecon.pdf

[og_spec]: http://bittorrent.org/beps/bep_0003.html
[other_spec]: https://wiki.theory.org/index.php/BitTorrentSpecification
[blackhat]: https://www.blackhat.com/presentations/bh-usa-09/BROOKS/BHUSA09-Brooks-BitTorrHacks-PAPER.pdf
