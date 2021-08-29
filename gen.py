events = [
    "Authenticated",
    "Pong",
    "Ready",
    "Message",
    "MessageUpdate",
    "MessageDelete",
    "ChannelCreate",
    "ChannelUpdate",
    "ChannelDelete",
    "ChannelGroupJoin",
    "ChannelGroupLeave",
    "ChannelStartTyping",
    "ChannelStopTyping",
    "ChannelAck",
    "ServerUpdate",
    "ServerDelete",
    "ServerMemberUpdate",
    "ServerMemberJoin",
    "ServerMemberLeave",
    "ServerRoleUpdate",
    "ServerRoleDelete",
    "UserUpdate",
    "UserRelationship",
]

print("switch messageType {")


parenthesis = "{}"
open = "{"
close = "}"

print('switch messageType {')

for event in events:
    print(f'case "{event}":')
    print(f'    o := {event}{parenthesis}')
    print(f'    err = json.Unmarshal(data, &o)')
    print(f'    if err != nil {open}')
    print(f'        return err')
    print(f'    {close}')
    print("")
    print(f'    rb.On{event}(o)')

print("default:")
print(f'    println(messageType + " not implemented")')
print(close)

print("")

for event in events:
    print(
        f'func (rb *RevoltBot) On{event}(o {event}) {parenthesis}')

# case "x":
#     o := x{}
#     rb.x(o)


# func (rb *RevoltBot) Onx(o x{}) {}
