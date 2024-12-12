# ----- Wormhole -----
# The following code creates a new TMUX session titled "wormhole", which as the
# name suggests, spins up a local development environment of Wormhole. This is
# powered by Tilt, and documentation around this setup can be found here:
# https://wormhole.com/docs/build/toolkit/tilt.
tmux kill-session -t wormhole
tmux new -d -s wormhole
tmux send-keys -t wormhole "cd core" ENTER
tmux send-keys -t wormhole "tilt down -- --wormchain" ENTER
tmux send-keys -t wormhole "tilt up -- --wormchain" ENTER

# ----- Noble -----
# The following code creates a new TMUX session title "noble", which as the
# name suggests, spins up a local development environment of Noble.
tmux kill-session -t noble
tmux new -d -s noble
tmux send-keys -t noble "cd .." ENTER
tmux send-keys -t noble "sh local.sh -r" ENTER
