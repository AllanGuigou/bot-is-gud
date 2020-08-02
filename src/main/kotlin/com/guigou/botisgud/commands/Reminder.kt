package com.guigou.botisgud.commands

import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.behavior.channel.createEmbed
import com.gitlab.kordlib.core.event.message.ReactionAddEvent
import com.gitlab.kordlib.core.on
import com.guigou.botisgud.models.Reminder
import com.guigou.botisgud.services.ReminderService
import com.guigou.botisgud.services.ReminderServiceImpl
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.launch
import kotlinx.coroutines.newSingleThreadContext
import java.util.*

class Reminder : Command {
    var reminderService: ReminderService = ReminderServiceImpl()

    override fun register(client: Kord) {

        client.on<ReactionAddEvent> {
            reminderService.add(Reminder(user.asUser().id, message.asMessage().content, "TODO", Date()))
        }

        val context = newSingleThreadContext("reminderServiceCollector")
        GlobalScope.launch(context) {
            reminderService.get().collect { value ->
                coroutineScope {
                    client.getUser(value.userId)!!.getDmChannel().createEmbed {
                        title = "Reminder"
                        description = value.message
                    }
                    reminderService.remove(value)
                }
            }
        }
    }
}
