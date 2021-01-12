package com.guigou.botisgud.commands

import com.gitlab.kordlib.common.entity.Snowflake
import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.behavior.channel.createEmbed
import com.gitlab.kordlib.core.event.message.ReactionAddEvent
import com.gitlab.kordlib.core.on
import com.gitlab.kordlib.rest.builder.message.EmbedBuilder
import com.guigou.botisgud.extensions.kord.link
import com.guigou.botisgud.extensions.logger
import com.guigou.botisgud.models.AbsoluteReminderTrigger
import com.guigou.botisgud.models.RelativeReminderTrigger
import com.guigou.botisgud.models.ReminderDto
import com.guigou.botisgud.models.ReminderTrigger
import com.guigou.botisgud.services.reminder.ReminderService
import kotlinx.coroutines.*
import kotlinx.coroutines.flow.collect
import me.tatarka.inject.annotations.Inject
import java.time.temporal.ChronoUnit

@Inject
class Reminder(private val service: ReminderService) : Command {
    companion object {
        val logger = logger()
    }

    private val reactions: Map<String, ReminderTrigger> = mapOf(
        Pair("⌚", RelativeReminderTrigger(1, ChronoUnit.HOURS)),
        Pair("☀️", AbsoluteReminderTrigger("0 9 * * *")),
        Pair("🌑", AbsoluteReminderTrigger("0 18 * * *"))
    )

    override suspend fun register(client: Kord, scope: CoroutineScope) {
        client.on<ReactionAddEvent> {
            val trigger = reactions[emoji.name] ?: return@on
            val dto = ReminderDto(userId, message.asMessage().content, link, trigger)
            service.add(dto)
        }

        backgroundWork(client, scope)
    }

    private fun backgroundWork(client: Kord, scope: CoroutineScope) {
        scope.launch(Dispatchers.Default) {
            service.get().collect {
                client.getUser(Snowflake(it.userId))!!.getDmChannel().createEmbed {
                    title = "Reminder"
                    description = it.message
                    fields.add(
                        EmbedBuilder.Field().apply {
                            this.name = "Original Message"
                            this.value = it.link
                        }
                    )
                }
                service.remove(it.id)
            }
        }
    }
}
