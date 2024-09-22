FROM scratch

COPY ./artifact/tw-stash /bin/tw-stash

ENTRYPOINT ["/bin/tw-stash"]
