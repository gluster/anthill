FROM docker.io/centos:7.5.1804
RUN yum update -y \
    && yum clean all

USER nobody

ADD build/_output/bin/anthill /usr/local/bin/anthill
