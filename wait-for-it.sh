#!/bin/bash
# wait-for-it.sh

# Standardized version - ensures consistent behavior and argument handling.

cmdname=$(basename $0)

TIMEOUT=15
QUIET=0
STRICT=0 # Default to not strict

echoerr() { if [[ $QUIET -ne 1 ]]; then echo "$@" 1>&2; fi }

usage()
{
    cat << USAGE >&2
Usage:
    $cmdname host:port [-s] [-t timeout] [-- command args]
    -h HOST | --host=HOST       Host or IP under test (ignored if host:port provided)
    -p PORT | --port=PORT       TCP port under test (ignored if host:port provided)
    -s | --strict               Only execute subcommand if the test succeeds
    -q | --quiet                Don't output any status messages
    -t TIMEOUT | --timeout=TIMEOUT
                                Timeout in seconds, zero for no timeout
    -- COMMAND ARGS             Execute command with args after the test finishes
USAGE
    exit 1
}

# process arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        *:* )
        HOSTPORT=(${1//:/ })
        HOST=${HOSTPORT[0]}
        PORT=${HOSTPORT[1]}
        if ! [[ $PORT =~ ^[0-9]+$ ]]; then
            echoerr "Error: Invalid port specification '$1'"
            usage
        fi
        shift 1
        ;;
        -s|--strict)
        STRICT=1
        shift 1
        ;;
        -q|--quiet)
        QUIET=1
        shift 1
        ;;
        -t)
        TIMEOUT="$2"
        if ! [[ $TIMEOUT =~ ^[0-9]+$ ]]; then
            echoerr "Error: Invalid timeout value '$TIMEOUT'"
            usage
        fi
        shift 2
        ;;
        --timeout=*)
        TIMEOUT="${1#*=}"
        if ! [[ $TIMEOUT =~ ^[0-9]+$ ]]; then
            echoerr "Error: Invalid timeout value '$TIMEOUT'"
            usage
        fi
        shift 1
        ;;
        --)
        shift
        CMD=("$@")
        break
        ;;
        --help)
        usage
        ;;
        *)
        echoerr "Unknown argument: $1"
        usage
        ;;
    esac
done

if [[ -z "$HOST" || -z "$PORT" ]]; then
    echoerr "Error: you need to provide a host and port to test."
    usage
fi

if [[ $TIMEOUT -lt 0 ]]; then
    echoerr "Error: timeout value must be 0 or greater"
    usage
fi

# Check dependencies. nc is preferred, fallback to bash /dev/tcp
WAIT_METHOD=""
# Check for nc (netcat-openbsd in Debian/Alpine)
if command -v nc >/dev/null 2>&1; then
    # Check if nc supports -z flag (OpenBSD variant does)
    if nc -z "$HOST" "$PORT" </dev/null >/dev/null 2>&1; then
        WAIT_METHOD="nc"
    # Add check for traditional netcat -w flag if needed, though -z is common
    # elif nc -w 1 "$HOST" "$PORT" < /dev/null > /dev/null 2>&1; then
    #    WAIT_METHOD="nc_traditional" # Requires different usage in wait_for
    fi
fi
# Fallback to bash /dev/tcp if nc is not suitable or available
if [[ -z "$WAIT_METHOD" ]] && [[ -n "$BASH_VERSION" ]]; then
    WAIT_METHOD="bash_dev_tcp"
fi

if [[ -z "$WAIT_METHOD" ]]; then
    echoerr "Error: this script requires nc (netcat) or bash with /dev/tcp support."
    exit 1
fi


wait_for()
{
    if [[ $TIMEOUT -gt 0 ]]; then
        echoerr "$cmdname: waiting $TIMEOUT seconds for $HOST:$PORT"
    else
        echoerr "$cmdname: waiting for $HOST:$PORT without timeout"
    fi
    start_ts=$(date +%s)
    while :
    do
        if [[ "$WAIT_METHOD" == "nc" ]]; then
            nc -z "$HOST" "$PORT" </dev/null >/dev/null 2>&1
            result=$?
        # Add handling for traditional nc if needed
        # elif [[ "$WAIT_METHOD" == "nc_traditional" ]]; then
        #     nc -w 1 "$HOST" "$PORT" < /dev/null > /dev/null 2>&1
        #     result=$?
        else # bash_dev_tcp
            (echo > /dev/tcp/$HOST/$PORT) >/dev/null 2>&1
            result=$?
        fi

        if [[ $result -eq 0 ]]; then
            end_ts=$(date +%s)
            echoerr "$cmdname: $HOST:$PORT is available after $((end_ts - start_ts)) seconds"
            break
        fi
        sleep 1
        if [[ $TIMEOUT -gt 0 ]]; then
            now_ts=$(date +%s)
            if [[ $((now_ts - start_ts)) -ge $TIMEOUT ]]; then
                echoerr "$cmdname: timeout occurred after waiting $TIMEOUT seconds for $HOST:$PORT"
                return 1 # Indicate timeout failure
            fi
        fi
    done
    return 0 # Indicate success
}

wait_for_wrapper()
{
    # In order to support SIGINT during timeout: http://unix.stackexchange.com/a/57692
    if [[ $QUIET -eq 1 ]]; then
        wait_for
        RESULT=$?
    else
        # Use timeout command if available (coreutils)
        # Note: timeout command might not be installed by default in slim/alpine
        # if command -v timeout >/dev/null 2>&1; then
        #      timeout $TIMEOUT bash -c wait_for # Requires bash -c to work with function
        #      RESULT=$?
        # else
            # Fallback to simple wait_for if timeout command is not available
            wait_for
            RESULT=$?
        # fi
    fi
    # Propagate the result of wait_for
    return $RESULT
}


wait_for_wrapper
RESULT=$?

# Execute command if provided
if [[ ${#CMD[@]} -gt 0 ]]; then
    # If wait failed and strict mode is enabled, exit
    if [[ $RESULT -ne 0 && $STRICT -eq 1 ]]; then
        echoerr "$cmdname: strict mode, refusing to execute subprocess"
        exit $RESULT
    fi
    # Use exec to replace the shell process with the command, passing args correctly
    exec "${CMD[@]}"
else
    # Exit with the result of the wait command if no command was specified
    exit $RESULT
fi

