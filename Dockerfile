FROM scratch

COPY --chmod=700 ./artifact/tw-stash /bin/tw-stash

ENTRYPOINT ["/bin/tw-stash"]
