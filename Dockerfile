FROM gsoci.azurecr.io/giantswarm/docker-kubectl:1.37.0

COPY crds/ /crds/

# The installer Job runs as an unprivileged user, so make the CRDs readable no matter what
# permissions they carry in the build context.
RUN chmod -R a+rX /crds
