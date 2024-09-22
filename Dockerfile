FROM scratch

COPY ./artifact/tw-stash /tw-stash

ENTRYPOINT ["/tw-stash"]
