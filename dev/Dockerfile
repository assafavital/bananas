FROM public.ecr.aws/u8k8d5f0/dev:ubuntu-latest

COPY dev/env /code/dev/env
COPY dev/cmds.sh /code/dev/cmds.sh
RUN apt-get update && \
    apt-get install cargo -y && \
    cargo install fblog && \
    export PATH=/root/.cargo/bin:$PATH
RUN printf 'source /code/dev/cmds.sh\nsource /code/dev/env\n' >> /root/.raftt_profile && \
    printf 'source /root/.raftt_profile\n' >> /root/.bash_profile && \
    printf 'source /root/.raftt_profile\n' >> /root/.bashrc
