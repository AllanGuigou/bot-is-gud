package com.guigou.botisgud.models

import java.time.Instant
import java.time.temporal.TemporalUnit

interface ReminderTrigger {
    fun timestamp(): Instant
}

class RelativeReminderTrigger(private val value: Long, private val unit: TemporalUnit) : ReminderTrigger {
    override fun timestamp(): Instant {
        return Instant.now()
                .plus(value, unit)
    }
}

class AbsoluteReminderTrigger() : ReminderTrigger {
    override fun timestamp(): Instant {
        // TODO: look into cron syntax
        return Instant.now()
    }
}
