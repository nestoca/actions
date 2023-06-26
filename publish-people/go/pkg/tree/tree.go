package tree

import (
	"github.com/nestoca/jac/pkg/live"
	"html"
	"regexp"
	"sort"
	"strings"
)

type Tree struct {
	Streams []*Stream
}

type Stream struct {
	*live.Group
	Name        string
	Description string
	Teams       []*Team
	Members     []*Member
}

type Team struct {
	*live.Group
	Name    string
	Members []*Member
}

type Member struct {
	*live.Person
	Name  string
	Email string
	Roles []*Role
}

type Role struct {
	*live.Group
	Name string
}

func NewTree(catalog *live.Catalog) *Tree {
	var streams []*Stream
	for _, streamGroup := range catalog.Root.Groups {
		// Only consider streams
		if streamGroup.Spec.Type != "stream" {
			continue
		}

		// Format description
		description := streamGroup.GetValueOrDefault("description", "")
		description = convertBulletPointsToHTMLList(description)
		description = convertLinks(description)

		// Add stream
		stream := &Stream{
			Group:       streamGroup,
			Name:        streamGroup.GetDisplayName(false, true),
			Description: description,
		}
		streams = append(streams, stream)

		// Determine people belonging directly to that stream
		for _, person := range streamGroup.Members {
			member := &Member{
				Person: person,
				Name:   person.GetDisplayName(false),
				Email:  person.Spec.Email,
			}
			stream.Members = append(stream.Members, member)

			// Determine person's roles
			for _, roleGroup := range person.Groups {
				if roleGroup.Spec.Type == "role" {
					member.Roles = append(member.Roles, &Role{
						Group: roleGroup,
						Name:  roleGroup.GetDisplayName(false, true),
					})
				}
			}
		}

		// Determine teams belonging to that stream
		sort.Slice(streamGroup.Children, func(i, j int) bool {
			return streamGroup.Children[i].Name < streamGroup.Children[j].Name
		})
		for _, teamGroup := range streamGroup.Children {
			// Only consider teams
			if teamGroup.Spec.Type != "team" {
				continue
			}

			// Add team
			team := &Team{
				Group: teamGroup,
				Name:  teamGroup.GetDisplayName(false, true),
			}
			stream.Teams = append(stream.Teams, team)

			for _, person := range teamGroup.Members {
				// Only consider members not already direct members of the stream
				if person.IsMemberOfGroup(streamGroup) {
					continue
				}

				member := Member{
					Person: person,
					Name:   person.GetDisplayName(false),
					Email:  person.Spec.Email,
				}
				team.Members = append(team.Members, &member)

				// Determine person's roles
				for _, group := range person.Groups {
					if group.Spec.Type == "role" {
						member.Roles = append(member.Roles, &Role{
							Group: group,
							Name:  group.GetDisplayName(false, true),
						})
					}
				}
			}

			sort.Slice(team.Members, func(i, j int) bool {
				return team.Members[i].Name < team.Members[j].Name
			})
		}
	}

	return &Tree{
		Streams: streams,
	}
}

func convertBulletPointsToHTMLList(bulletPoints string) string {
	if bulletPoints == "" {
		return ""
	}
	bulletPoints = html.EscapeString(bulletPoints)
	lines := strings.Split(bulletPoints, "\n")
	result := "<ul>\n"
	result += convertLinesToHTML(lines, 0)
	result += "</ul>"
	return result
}

func convertLinesToHTML(lines []string, indentLevel int) string {
	html := ""
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmedLine := strings.TrimRight(line, " ")
		if trimmedLine == "" {
			continue
		}

		indent := strings.Repeat("  ", indentLevel)
		item := strings.TrimPrefix(trimmedLine, "- ")
		html += indent + "<li>" + item

		nextIndentLevel := indentLevel + 1
		childLines := getChildLines(lines[i+1:], nextIndentLevel)
		if len(childLines) > 0 {
			html += "<ul>\n"
			html += convertLinesToHTML(childLines, nextIndentLevel)
			html += indent + "  </ul>"
			i += len(childLines)
		}

		html += "</li>\n"
	}
	return html
}

func getChildLines(lines []string, indentLevel int) []string {
	childLines := make([]string, 0)
	for _, line := range lines {
		trimmedLine := strings.TrimRight(line, " ")
		if strings.HasPrefix(trimmedLine, strings.Repeat("  ", indentLevel)+"- ") {
			childLines = append(childLines, strings.TrimPrefix(trimmedLine, strings.Repeat("  ", indentLevel)+"- "))
		} else if trimmedLine != "" {
			break
		}
	}
	return childLines
}

var linkRegex = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)

func convertLinks(input string) string {
	return linkRegex.ReplaceAllString(input, `<a href="$2">$1</a>`)
}
