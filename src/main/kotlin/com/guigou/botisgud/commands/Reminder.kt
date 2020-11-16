package com.guigou.botisgud.commands

import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.behavior.channel.createEmbed
import com.gitlab.kordlib.core.event.message.ReactionAddEvent
import com.gitlab.kordlib.core.on
import com.gitlab.kordlib.rest.builder.message.EmbedBuilder
import com.guigou.botisgud.extensions.kord.link
import com.guigou.botisgud.models.AbsoluteReminderTrigger
import com.guigou.botisgud.models.RelativeReminderTrigger
import com.guigou.botisgud.models.ReminderDto
import com.guigou.botisgud.models.ReminderTrigger
import com.guigou.botisgud.services.ReminderService
import com.guigou.botisgud.services.ReminderServiceImpl
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.launch
import kotlinx.coroutines.newSingleThreadContext
import java.time.temporal.ChronoUnit

class Reminder(private val service: ReminderService = ReminderServiceImpl()) : Command {
    private val reactions: Map<String, ReminderTrigger> = mapOf(
        Pair("⌚", RelativeReminderTrigger(1, ChronoUnit.HOURS)),
        Pair("☀️", AbsoluteReminderTrigger()),
        Pair("🌑", AbsoluteReminderTrigger())
    )

    override fun register(client: Kord) {
        client.on<ReactionAddEvent> {
            if (!reactions.containsKey(emoji.name)) {
                return@on
            }

            val trigger = reactions[emoji.name]
            val reminderDto = ReminderDto(userId, message.asMessage().content, link)
            service.add(reminderDto, trigger!!)
        }

        val context = newSingleThreadContext("reminderServiceCollector")
        GlobalScope.launch(context) {
            service.get().collect { value ->
                coroutineScope {
                    client.getUser(value.userId)!!.getDmChannel().createEmbed {
                        title = "Reminder"
                        description = value.message
                        fields.add(
                            EmbedBuilder.Field().apply {
                                this.name = "Original Message"
                                this.value = value.link.toString()
                            }
                        )
                    }
                    service.remove(value)
                }
            }
        }
    }
}
