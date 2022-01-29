using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using Microsoft.EntityFrameworkCore;

public class PresenceContext : DbContext
{
	public DbSet<PresenceEvent> PresenceEvents { get; set; }

	public string Path { get; }

	public PresenceContext()
	{
		Path = "presence.db";
	}

	protected override void OnConfiguring(DbContextOptionsBuilder options) =>
		options.UseSqlite($"Data Source={Path}");
}

public class PresenceEvent
{
	[Key]
	public int Id { get; set; }
	public DateTime Timestamp { get; set; }
	public long UserId { get; set; }
	public string Status { get; set; }
}