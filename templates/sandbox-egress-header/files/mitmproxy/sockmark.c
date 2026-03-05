/*
 * SO_MARK LD_PRELOAD library for mitmproxy loop prevention.
 *
 * Hooks socket() to set SO_MARK=1 on all new IPv4/IPv6 sockets.
 * When loaded via LD_PRELOAD into mitmproxy, this marks all of
 * mitmproxy's outbound packets so iptables can skip them
 * (preventing redirect loops in transparent proxy mode).
 *
 * Requires CAP_NET_ADMIN (run mitmproxy as root in the sandbox).
 *
 * Build: gcc -shared -fPIC -o sockmark.so sockmark.c -ldl
 */
#define _GNU_SOURCE
#include <sys/socket.h>
#include <dlfcn.h>

int socket(int domain, int type, int protocol) {
    int (*real_socket)(int, int, int) = dlsym(RTLD_NEXT, "socket");
    int fd = real_socket(domain, type, protocol);
    if (fd >= 0 && (domain == AF_INET || domain == AF_INET6)) {
        int mark = 1;
        setsockopt(fd, SOL_SOCKET, SO_MARK, &mark, sizeof(mark));
    }
    return fd;
}
