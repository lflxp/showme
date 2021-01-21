// http://dockone.io/article/109
// https://docs.docker.com/engine/api/v1.24/#2-errors
// https://godoc.org/github.com/docker/docker/client#FromEnv
package utils

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type docker struct {
	host    string
	port    string
	version string
}

func NewDockerCLI(host, port, version string) *docker {
	if version == "" {
		version = "v1.24"
	}

	return &docker{
		host:    host,
		port:    port,
		version: version,
	}
}

// Query parameters:
// all – 1/True/true or 0/False/false, Show all containers. Only running containers are shown by default (i.e., this defaults to false)
// limit – Show limit last created containers, include non-running ones.
// since – Show only containers created since Id, include non-running ones.
// before – Show only containers created before Id, include non-running ones.
// size – 1/True/true or 0/False/false, Show the containers sizes
// filters - a JSON encoded value of the filters (a map[string][]string) to process on the containers list. Available filters:
// exited=<int>; -- containers with exit code of <int> ;
// status=(created	restarting	running	paused	exited	dead)
// label=key or label="key=value" of a container label
// isolation=(default	process	hyperv) (Windows daemon only)
// ancestor=(<image-name>[:<tag>], <image id> or <image@digest>)
// before=(<container id> or <container name>)
// since=(<container id> or <container name>)
// volume=(<volume name> or <mount point destination>)
// network=(<network id> or <network name>)
// Status codes:
// 200 – no error
// 400 – bad parameter
// 500 – server error
func (d *docker) ListContainers() (result interface{}, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/json?all=1", d.host, d.port, d.version)
	log.Debug(url)
	data, err := Get(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}

	return
}

// Create a container
// Hostname - A string value containing the hostname to use for the container. This must be a valid RFC 1123 hostname.
// Domainname - A string value containing the domain name to use for the container.
// User - A string value specifying the user inside the container.
// AttachStdin - Boolean value, attaches to stdin.
// AttachStdout - Boolean value, attaches to stdout.
// AttachStderr - Boolean value, attaches to stderr.
// Tty - Boolean value, Attach standard streams to a tty, including stdin if it is not closed.
// OpenStdin - Boolean value, opens stdin,
// StdinOnce - Boolean value, close stdin after the 1 attached client disconnects.
// Env - A list of environment variables in the form of ["VAR=value", ...]
// Labels - Adds a map of labels to a container. To specify a map: {"key":"value", ... }
// Cmd - Command to run specified as a string or an array of strings.
// Entrypoint - Set the entry point for the container as a string or an array of strings.
// Image - A string specifying the image name to use for the container.
// Volumes - An object mapping mount point paths (strings) inside the container to empty objects.
// Healthcheck - A test to perform to check that the container is healthy.
// Test - The test to perform. Possible values are: + {} inherit healthcheck from image or parent image + {"NONE"} disable healthcheck + {"CMD", args...} exec arguments directly + {"CMD-SHELL", command} run command with system’s default shell
// Interval - The time to wait between checks in nanoseconds. It should be 0 or at least 1000000 (1 ms). 0 means inherit.
// Timeout - The time to wait before considering the check to have hung. It should be 0 or at least 1000000 (1 ms). 0 means inherit.
// Retries - The number of consecutive failures needed to consider a container as unhealthy. 0 means inherit.
// StartPeriod - The time to wait for container initialization before starting health-retries countdown in nanoseconds. It should be 0 or at least 1000000 (1 ms). 0 means inherit.
// WorkingDir - A string specifying the working directory for commands to run in.
// NetworkDisabled - Boolean value, when true disables networking for the container
// ExposedPorts - An object mapping ports to an empty object in the form of: "ExposedPorts": { "<port>/<tcp|udp>: {}" }
// StopSignal - Signal to stop a container as a string or unsigned integer. SIGTERM by default.
// HostConfig
// Binds – A list of volume bindings for this container. Each volume binding is a string in one of these forms:
// host-src:container-dest to bind-mount a host path into the container. Both host-src, and container-dest must be an absolute path.
// host-src:container-dest:ro to make the bind mount read-only inside the container. Both host-src, and container-dest must be an absolute path.
// volume-name:container-dest to bind-mount a volume managed by a volume driver into the container. container-dest must be an absolute path.
// volume-name:container-dest:ro to mount the volume read-only inside the container. container-dest must be an absolute path.
// Tmpfs – A map of container directories which should be replaced by tmpfs mounts, and their corresponding mount options. A JSON object in the form { "/run": "rw,noexec,nosuid,size=65536k" }.
// Links - A list of links for the container. Each link entry should be in the form of container_name:alias.
// Memory - Memory limit in bytes.
// MemorySwap - Total memory limit (memory + swap); set -1 to enable unlimited swap. You must use this with memory and make the swap value larger than memory.
// MemoryReservation - Memory soft limit in bytes.
// KernelMemory - Kernel memory limit in bytes.
// CpuPercent - An integer value containing the usable percentage of the available CPUs. (Windows daemon only)
// CpuShares - An integer value containing the container’s CPU Shares (ie. the relative weight vs other containers).
// CpuPeriod - The length of a CPU period in microseconds.
// CpuQuota - Microseconds of CPU time that the container can get in a CPU period.
// CpusetCpus - String value containing the cgroups CpusetCpus to use.
// CpusetMems - Memory nodes (MEMs) in which to allow execution (0-3, 0,1). Only effective on NUMA systems.
// IOMaximumBandwidth - Maximum IO absolute rate in terms of IOps.
// IOMaximumIOps - Maximum IO absolute rate in terms of bytes per second.
// BlkioWeight - Block IO weight (relative weight) accepts a weight value between 10 and 1000.
// BlkioWeightDevice - Block IO weight (relative device weight) in the form of: "BlkioWeightDevice": [{"Path": "device_path", "Weight": weight}]
// BlkioDeviceReadBps - Limit read rate (bytes per second) from a device in the form of: "BlkioDeviceReadBps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceReadBps": [{"Path": "/dev/sda", "Rate": "1024"}]"
// BlkioDeviceWriteBps - Limit write rate (bytes per second) to a device in the form of: "BlkioDeviceWriteBps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceWriteBps": [{"Path": "/dev/sda", "Rate": "1024"}]"
// BlkioDeviceReadIOps - Limit read rate (IO per second) from a device in the form of: "BlkioDeviceReadIOps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceReadIOps": [{"Path": "/dev/sda", "Rate": "1000"}]
// BlkioDeviceWriteIOps - Limit write rate (IO per second) to a device in the form of: "BlkioDeviceWriteIOps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceWriteIOps": [{"Path": "/dev/sda", "Rate": "1000"}]
// MemorySwappiness - Tune a container’s memory swappiness behavior. Accepts an integer between 0 and 100.
// OomKillDisable - Boolean value, whether to disable OOM Killer for the container or not.
// OomScoreAdj - An integer value containing the score given to the container in order to tune OOM killer preferences.
// PidMode - Set the PID (Process) Namespace mode for the container; "container:<name|id>": joins another container’s PID namespace "host": use the host’s PID namespace inside the container
// PidsLimit - Tune a container’s pids limit. Set -1 for unlimited.
// PortBindings - A map of exposed container ports and the host port they should map to. A JSON object in the form { <port>/<protocol>: [{ "HostPort": "<port>" }] } Take note that port is specified as a string and not an integer value.
// PublishAllPorts - Allocates an ephemeral host port for all of a container’s exposed ports. Specified as a boolean value.
// Ports are de-allocated when the container stops and allocated when the container starts. The allocated port might be changed when restarting the container.
// The port is selected from the ephemeral port range that depends on the kernel. For example, on Linux the range is defined by /proc/sys/net/ipv4/ip_local_port_range.
// Privileged - Gives the container full access to the host. Specified as a boolean value.
// ReadonlyRootfs - Mount the container’s root filesystem as read only. Specified as a boolean value.
// Dns - A list of DNS servers for the container to use.
// DnsOptions - A list of DNS options
// DnsSearch - A list of DNS search domains
// ExtraHosts - A list of hostnames/IP mappings to add to the container’s /etc/hosts file. Specified in the form ["hostname:IP"].
// VolumesFrom - A list of volumes to inherit from another container. Specified in the form <container name>[:<ro|rw>]
// CapAdd - A list of kernel capabilities to add to the container.
// Capdrop - A list of kernel capabilities to drop from the container.
// GroupAdd - A list of additional groups that the container process will run as
// RestartPolicy – The behavior to apply when the container exits. The value is an object with a Name property of either "always" to always restart, "unless-stopped" to restart always except when user has manually stopped the container or "on-failure" to restart only when the container exit code is non-zero. If on-failure is used, MaximumRetryCount controls the number of times to retry before giving up. The default is not to restart. (optional) An ever increasing delay (double the previous delay, starting at 100mS) is added before each restart to prevent flooding the server.
// UsernsMode - Sets the usernamespace mode for the container when usernamespace remapping option is enabled. supported values are: host.
// NetworkMode - Sets the networking mode for the container. Supported standard values are: bridge, host, none, and container:<name|id>. Any other value is taken as a custom network’s name to which this container should connect to.
// Devices - A list of devices to add to the container specified as a JSON object in the form { "PathOnHost": "/dev/deviceName", "PathInContainer": "/dev/deviceName", "CgroupPermissions": "mrw"}
// Ulimits - A list of ulimits to set in the container, specified as { "Name": <name>, "Soft": <soft limit>, "Hard": <hard limit> }, for example: Ulimits: { "Name": "nofile", "Soft": 1024, "Hard": 2048 }
// Sysctls - A list of kernel parameters (sysctls) to set in the container, specified as { <name>: <Value> }, for example: { "net.ipv4.ip_forward": "1" }
// SecurityOpt: A list of string values to customize labels for MLS systems, such as SELinux.
// StorageOpt: Storage driver options per container. Options can be passed in the form {"size":"120G"}
// LogConfig - Log configuration for the container, specified as a JSON object in the form { "Type": "<driver_name>", "Config": {"key1": "val1"}}. Available types: json-file, syslog, journald, gelf, fluentd, awslogs, splunk, etwlogs, none. json-file logging driver.
// CgroupParent - Path to cgroups under which the container’s cgroup is created. If the path is not absolute, the path is considered to be relative to the cgroups path of the init process. Cgroups are created if they do not already exist.
// VolumeDriver - Driver that this container users to mount volumes.
// ShmSize - Size of /dev/shm in bytes. The size must be greater than 0. If omitted the system uses 64MB.
// Query parameters:
// name – Assign the specified name to the container. Must match /?[a-zA-Z0-9_-]+.
// Status codes:
// 201 – no error
// 400 – bad parameter
// 404 – no such container
// 406 – impossible to attach (container not running)
// 409 – conflict
// 500 – server error
func (d *docker) CreateContainer(data interface{}) (result []byte, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/create", d.host, d.port, d.version)
	log.Debug(url)
	result, err = Post(url, data, "application/json")
	return
}

// 监控容器，使用容器id获取该容器的底层信息
func (d *docker) ContainerJson(id string) (result interface{}, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/%s/json", d.host, d.port, d.version, id)
	log.Debug(url)
	data, err := Get(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}

	return
}

// 进程列表。获取容器内进程的清单：
func (d *docker) ContainerTop(id string) (result interface{}, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/%s/top", d.host, d.port, d.version, id)
	log.Debug(url)
	data, err := Get(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}

	return
}

// 容器日志。获取容器的标准输出和错误日志
// Query parameters:
// details - 1/True/true or 0/False/false, Show extra details provided to logs. Default false.
// follow – 1/True/true or 0/False/false, return stream. Default false.
// stdout – 1/True/true or 0/False/false, show stdout log. Default false.
// stderr – 1/True/true or 0/False/false, show stderr log. Default false.
// since – UNIX timestamp (integer) to filter logs. Specifying a timestamp will only output log-entries since that timestamp. Default: 0 (unfiltered)
// timestamps – 1/True/true or 0/False/false, print timestamps for every log line. Default false.
// tail – Output specified number of lines at the end of logs: all or <number>. Default all.
// Status codes:
// 101 – no error, hints proxy about hijacking
// 200 – no error, no upgrade header found
// 404 – no such container
// 500 – server error
func (d *docker) ContainerLogs(id string) (result interface{}, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/%s/logs?stderr=1&stdout=1&timestamps=1&follow=0", d.host, d.port, d.version, id)
	log.Debug(url)
	data, err := Get(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}

	return
}

// 启动容器
// Query parameters:
// detachKeys – Override the key sequence for detaching a container. Format is a single character [a-Z] or ctrl-<value> where <value> is one of: a-z, @, ^, [, , or _.
// Status codes:
// 204 – no error
// 304 – container already started
// 404 – no such container
// 500 – server error
func (d *docker) ContainerStart(id string, data interface{}) (result []byte, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/%s/start", d.host, d.port, d.version, id)
	log.Debug(url)
	result, err = Post(url, data, "application/json")
	return
}

// 停止容器
// Query parameters:
// t – number of seconds to wait before killing the container
// Status codes:
// 204 – no error
// 304 – container already stopped
// 404 – no such container
// 500 – server error
func (d *docker) ContainerStop(id string) (result []byte, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/%s/stop?t=5", d.host, d.port, d.version, id)
	log.Debug(url)
	result, err = Post(url, nil, "application/json")
	return
}

// 重启容器
// Query parameters:
// t – number of seconds to wait before killing the container
// Status codes:
// 204 – no error
// 404 – no such container
// 500 – server error
func (d *docker) ContainerRestart(id string) (result []byte, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/%s/restart?t=5", d.host, d.port, d.version, id)
	log.Debug(url)
	result, err = Post(url, nil, "application/json")
	return
}

// 终止容器
// Query parameters:
// signal - Signal to send to the container: integer or string like SIGINT. When not set, SIGKILL is assumed and the call waits for the container to exit.
// Status codes:
// 204 – no error
// 404 – no such container
// 500 – server error
func (d *docker) ContainerKill(id string) (result []byte, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/%s/kill", d.host, d.port, d.version, id)
	log.Debug(url)
	result, err = Post(url, nil, "application/json")
	return
}

// TODO: update
