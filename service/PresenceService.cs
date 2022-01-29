
using Grpc.Core;
using Presence;
class PresenceServiceImpl : PresenceService.PresenceServiceBase
{
	public PresenceServiceImpl()
	{
	}

	public override async Task<EventResponse> TrackEvent(EventRequest request, ServerCallContext context)
	{
		using var db = new PresenceContext();

		try
		{
			Console.WriteLine($"{request.Timestamp} {request.User} {request.Status}");
			await db.PresenceEvents.AddAsync(new PresenceEvent
			{
				Timestamp = request.Timestamp.ToDateTime(),
				UserId = long.Parse(request.User),
				Status = request.Status,
			});
			await db.SaveChangesAsync();

			Console.WriteLine(db.PresenceEvents.Count());

			return new EventResponse();
		}
		catch (Exception ex)
		{
			Console.WriteLine(ex);
			throw ex;
		}
	}
}