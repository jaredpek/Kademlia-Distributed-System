FROM larjim/kademlialab

RUN apt-get update -y
RUN apt-get install -y iputils-ping
RUN apt-get install -y arp-scan
# arp-scan -l

# Remove old version of Go
RUN rm -rf /usr/local/go

# Download and install Go
RUN wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz && \
    rm go1.22.1.linux-amd64.tar.gz

# Set the PATH environment variable
ENV PATH=$PATH:/usr/local/go/bin

RUN go version

ADD main.go .
ADD go.mod .
ADD /kademlia/ ./kademlia/

#CMD go run /home/main.go or something similar to run by default i believe

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab
