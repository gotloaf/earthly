
script_full_path=$(dirname "$0")

$script_full_path/build_docs.sh
$script_full_path/build_server.sh
./build/server
