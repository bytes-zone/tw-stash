FROM scratch

COPY --chmod=700 ./artifact/tw-stash /tw-stash

ENTRYPOINT ["/tw-stash"]
