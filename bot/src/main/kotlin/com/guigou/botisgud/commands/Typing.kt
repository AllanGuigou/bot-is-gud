package com.guigou.botisgud.commands

import com.guigou.botisgud.extensions.logger
import dev.kord.core.Kord
import dev.kord.core.behavior.channel.MessageChannelBehavior
import dev.kord.core.entity.User
import dev.kord.core.event.channel.TypingStartEvent
import dev.kord.core.event.message.MessageCreateEvent
import dev.kord.core.on
import kotlinx.coroutines.CoroutineScope
import java.time.Instant
import kotlin.random.Random

class Typing : Command {
    companion object {
        var logger = logger()
    }

    private var triggered: Instant = Instant.MIN;

    private suspend fun type(user: User?, channel: MessageChannelBehavior) {
        if (user == null || user.asUser().isBot) {
            return
        }

        if (triggered.isAfter(Instant.now().minusSeconds(60))) {
            return
        }

        val chance = Random.nextInt(0, 100)
        if (chance > 20) {
            return
        }

        channel.type()
        triggered = Instant.now()
    }

    override suspend fun register(client: Kord, scope: CoroutineScope) {
        client.on<MessageCreateEvent> {
            logger.info("event: MESSAGE_CREATE user: ${member?.id}")
            type(member, message.channel)
        }
        client.on<TypingStartEvent> {
            logger.info("event: TYPING_START user: ${user.id}")
            type(user.asUser(), channel)
        }
    }
}
