package values

import (
	"github.com/nestoca/jac/pkg/live"
)

type Values struct {
	Streams []*Stream
}

type Stream struct {
	Group   *live.Group
	Teams   []*Team
	Members []*Member
}

type Team struct {
	Group   *live.Group
	Members []*Member
}

type Member struct {
	Person *live.Person
	Roles  []*live.Group
}

func NewValues(catalog *live.Catalog) *Values {
	var streams []*Stream
	for _, group := range catalog.Root.Groups {
		// Only consider streams
		if group.Spec.Type != "stream" {
			continue
		}

		// Add stream
		stream := &Stream{
			Group: group,
		}
		streams = append(streams, stream)

		// Determine people belonging directly to that stream
		for _, person := range group.Members {
			member := &Member{
				Person: person,
			}
			stream.Members = append(stream.Members, member)

			// Determine person's roles
			for _, group := range person.Groups {
				if group.Spec.Type == "role" {
					member.Roles = append(member.Roles, group)
				}
			}
		}

		// Determine teams belonging to that stream
		for _, group := range group.Children {
			// Only consider teams
			if group.Spec.Type != "team" {
				continue
			}

			// Add team
			team := &Team{
				Group: group,
			}
			stream.Teams = append(stream.Teams, team)

			for _, person := range group.Members {
				// Only consider members not already direct members of the stream
				if person.IsMemberOfGroup(stream.Group) {
					continue
				}

				member := Member{
					Person: person,
				}
				team.Members = append(team.Members, &member)

				// Determine person's roles
				for _, group := range person.Groups {
					if group.Spec.Type == "role" {
						member.Roles = append(member.Roles, group)
					}
				}
			}
		}
	}

	return &Values{
		Streams: streams,
	}
}
