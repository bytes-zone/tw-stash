FROM cgr.dev/chainguard/glibc-dynamic

COPY --chown=nonroot:nonroot ./artifact/tw-stash /usr/local/bin/tw-stash

ENTRYPOINT ["/usr/local/bin/tw-stash"]
