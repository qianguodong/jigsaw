ARG base_image
ARG debug_image

FROM $base_image
ARG USER_UID=10001
COPY . /go/bin/
ENTRYPOINT ["/go/bin/designer", "--config=designer.yaml"]
USER ${USER_UID}

FROM $debug_image AS debug
ARG USER_UID=10001
COPY . /go/bin/
ENTRYPOINT ["/go/bin/dlv", "exec", "/go/bin/designer", "--config=designer.yaml", "--headless", "--listen=:12345", "--api-version=2", "--accept-multiclient", "--log", "--"]
USER ${USER_UID}
