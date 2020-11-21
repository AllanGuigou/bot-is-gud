package com.guigou.botisgud.commands

import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.event.channel.TypingStartEvent
import com.gitlab.kordlib.core.on
import java.time.Instant
import kotlin.random.Random

class Typing : Command {
    private var triggered: Instant = Instant.MIN;

    override fun register(client: Kord) {
        client.on<TypingStartEvent> {
            if (user.asUser().isBot == true) {
                return@on
            }

            if (triggered.isAfter(Instant.now().minusSeconds(60))) {
                return@on
            }

            if (Random.nextInt(0, 100) < 20) {
                return@on
            }

            channel.type()
            triggered = Instant.now()
        }
    }
}
