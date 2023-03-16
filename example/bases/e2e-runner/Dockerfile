FROM alpine:latest

# install kubectl
RUN apk add docker curl bash && \
  curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.23.5/bin/linux/amd64/kubectl && \
  chmod +x ./kubectl && \
  mv ./kubectl /usr/local/bin/kubectl && \
  kubectl version --client

# install GRGate
RUN curl -LO https://github.com/FikaWorks/grgate/releases/download/v0.6.2/grgate_0.6.2_linux_amd64.tar.gz && \
  tar -xvf grgate_0.6.2_linux_amd64.tar.gz && \
  chmod +x ./grgate && \
  mv ./grgate /usr/local/bin/grgate && \
  rm -rf grgate_0.6.2_linux_amd64.tar.gz && \
  grgate version

COPY e2e-runner.sh /e2e-runner.sh

CMD ["/e2e-runner.sh"]
