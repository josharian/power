usage: `nobatt cmd [args]`

nobatt runs cmd.

It suspends cmd when it detects that the laptop is running on battery power,
and resumes it when it detects that the laptop is using wall power.
If you want to override nobatt, send SIGUSR1 to force cmd to resume,
or SIGUSR2 to force cmd to suspend.

Note that nobatt only controls cmd, not any other processes that cmd starts.
