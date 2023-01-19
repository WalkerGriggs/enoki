<h1 align="center">Enoki</h1>
<p align=center>
  <img alt="Enoki mushrooms, illustrated" width="250" height="250" src="./assets/enoki.jpeg">
</p>

## Goals and Philosophy

Too many media servers prioritize features and style above unambiguously serving your media. Some go as far as gathering your identifying information and locking core-functionality behind pay walls.

Not Enoki.

Enoki is written with three goals in mind:
1. total ownership of your content and metadata
2. streaming performance no matter the resource requirements
3. free without exception, _liber et gratis_ 

## Features and Design

- Playback security
  - Signed URLS

    All manifest and segment requests must be signed with a valid JWT. The JWT is generated when a client authenticates and requests the initial playback URL. That playback URL is pre-signed and, when requested, returns the manifest file.

    Each segment in that manifest file has been signed with that same token.

    The JWT contains client metadata like originating IP, media data like the asset ID, and security fields like an expiration and the signing key ID. These fields gives enoki a rich request history, and we can revoked specific targets if needed.

  - HLSe

    All segments, by default, are encryped with AES-128. The key is specified in the manifest and is rotated every 10 minutes.

- Rendition ladders and ABR

  Enoki prioritizes performance over storage space. True JIT across an arbitrary number of streams is difficult, especially in home application where we're severly limited by the number of encoders of CPU cores.

  By default, Enoki supports static rendition latters for standard, HD resolutions. We'd like to support per-title encoding using the convex hull of PSNR over bitrate.

- Storage classes

    Enoki does not store media or manage the local file system directly*. It implements a number of `StorageClasses` like S3, NFS, SMB, etc. Mezzanine files and transcoded segments are written to and streamed from those respective storage classes. Any class which does not support adequate access models are proxied through Enoki.

    \* Enoki does implement a local StorageClass, but this is used for development purposes only and should not be used into an actual deployment.

