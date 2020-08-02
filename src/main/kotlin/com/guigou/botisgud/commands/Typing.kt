package com.guigou.botisgud.commands

import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.event.channel.TypingStartEvent
import com.gitlab.kordlib.core.on

class Typing : Command {
    override fun register(client: Kord) {
        client.on<TypingStartEvent> {
            if (user.asUser().isBot == true) {
                return@on
            }

            channel.type()
        }
    }
}
