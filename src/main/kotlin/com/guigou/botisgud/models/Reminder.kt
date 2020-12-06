package com.guigou.botisgud.models

import java.time.Instant

// TODO: can Snowflake be serialized? if so then convert userId back to Snowflake.
data class Reminder(
    val id: Int,
    val userId: Long,
    val message: String,
    val link: String,
    val timestamp: Instant
) {
}
