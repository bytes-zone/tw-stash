FROM scratch

COPY ./target/release/tw-stash /tw-stash

ENTRYPOINT ["/tw-stash"]
