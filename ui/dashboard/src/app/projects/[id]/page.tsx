'use client';

import { useEffect, useState, useCallback } from 'react';
import { useParams } from 'next/navigation';
import {
  Badge, Button, Card, Checkbox, Grid, Group, Loader, Menu, NumberInput,
  Progress, ScrollArea, Stack, Text, Timeline, Title,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import {
  IconBulb, IconCheck, IconChevronDown, IconDatabase, IconEdit,
  IconPlayerPlay, IconSearch, IconSettings, IconX,
} from '@tabler/icons-react';
import Link from 'next/link';
import Shell from '@/components/layout/AppShell';
import { api, DiscoveryResult, DiscoveryRunStatus, Project } from '@/lib/api';

export default function ProjectPage() {
  const { id } = useParams<{ id: string }>();
  const [project, setProject] = useState<Project | null>(null);
  const [discoveries, setDiscoveries] = useState<DiscoveryResult[]>([]);
  const [run, setRun] = useState<DiscoveryRunStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [triggering, setTriggering] = useState(false);
  const [analysisAreas, setAnalysisAreas] = useState<{ id: string; name: string }[]>([]);
  const [selectedAreas, setSelectedAreas] = useState<string[]>([]);
  const [maxSteps, setMaxSteps] = useState(100);

  useEffect(() => {
    Promise.all([
      api.getProject(id).then((p) => {
        setProject(p);
        return api.getAnalysisAreas(p.domain, p.category)
          .then((areas) => setAnalysisAreas(areas.map((a) => ({ id: a.id, name: a.name }))));
      }),
      api.listDiscoveries(id).then(setDiscoveries).catch(() => []),
    ])
      .catch((e) => notifications.show({ title: 'Error', message: e.message, color: 'red' }))
      .finally(() => setLoading(false));
  }, [id]);

  const pollStatus = useCallback(async () => {
    try {
      const status = await api.getProjectStatus(id);
      if (status?.run) setRun(status.run as unknown as DiscoveryRunStatus);
    } catch { /* ignore */ }
  }, [id]);

  useEffect(() => {
    if (!run || (run.status !== 'running' && run.status !== 'pending')) return;
    const interval = setInterval(pollStatus, 2000);
    return () => clearInterval(interval);
  }, [run, pollStatus]);

  useEffect(() => { pollStatus(); }, [pollStatus]);

  const handleTrigger = async (areas?: string[]) => {
    setTriggering(true);
    try {
      const opts: { areas?: string[]; max_steps?: number } = {};
      if (areas && areas.length > 0) opts.areas = areas;
      if (maxSteps !== 100) opts.max_steps = maxSteps;

      const result = await api.triggerDiscovery(id, Object.keys(opts).length > 0 ? opts : undefined);
      if (result.run_id) {
        const newRun = await api.getRun(result.run_id);
        setRun(newRun);
      }
      notifications.show({ title: 'Discovery started', message: `${maxSteps} steps`, color: 'blue' });
    } catch (e: unknown) {
      notifications.show({ title: 'Error', message: (e as Error).message, color: 'red' });
    } finally {
      setTriggering(false);
      setSelectedAreas([]);
    }
  };

  if (loading) return <Shell><Loader /></Shell>;
  if (!project) return <Shell><Text>Project not found</Text></Shell>;

  const isRunning = run && (run.status === 'running' || run.status === 'pending');
  const latestDiscovery = discoveries.length > 0 ? discoveries[0] : null;

  return (
    <Shell>
      <Stack gap="lg">
        {/* Header */}
        <Group justify="space-between">
          <div>
            <Title order={2}>{project.name}</Title>
            <Group gap="xs" mt={4}>
              <Badge variant="light">{project.domain}</Badge>
              <Badge variant="light" color="blue">{project.category}</Badge>
              {project.description && <Text size="xs" c="dimmed">{project.description}</Text>}
            </Group>
          </div>
          <Group>
            <Button variant="subtle" component={Link} href={`/projects/${id}/prompts`}
              leftSection={<IconEdit size={16} />} size="sm">Prompts</Button>
            <Button variant="subtle" component={Link} href={`/projects/${id}/settings`}
              leftSection={<IconSettings size={16} />} size="sm">Settings</Button>

            <Menu shadow="md" width={280} disabled={!!isRunning}>
              <Menu.Target>
                <Button leftSection={<IconPlayerPlay size={16} />}
                  rightSection={<IconChevronDown size={14} />}
                  loading={triggering} disabled={!!isRunning}>
                  {isRunning ? 'Running...' : 'Run Discovery'}
                </Button>
              </Menu.Target>
              <Menu.Dropdown>
                <Menu.Label>Exploration steps</Menu.Label>
                <div style={{ padding: '4px 12px 8px' }}>
                  <NumberInput size="xs" value={maxSteps} onChange={(v) => setMaxSteps(Number(v) || 100)}
                    min={5} max={500} step={5} description="More steps = more comprehensive" />
                </div>
                <Menu.Divider />
                <Menu.Item onClick={() => handleTrigger()}>Run All Areas</Menu.Item>
                <Menu.Divider />
                <Menu.Label>Select areas</Menu.Label>
                {analysisAreas.map((area) => (
                  <Menu.Item key={area.id} closeMenuOnClick={false}>
                    <Checkbox label={area.name} checked={selectedAreas.includes(area.id)}
                      onChange={(e) => {
                        if (e.currentTarget.checked) setSelectedAreas([...selectedAreas, area.id]);
                        else setSelectedAreas(selectedAreas.filter((a) => a !== area.id));
                      }} />
                  </Menu.Item>
                ))}
                {selectedAreas.length > 0 && (
                  <>
                    <Menu.Divider />
                    <Menu.Item color="blue" onClick={() => handleTrigger(selectedAreas)}>
                      Run Selected ({selectedAreas.length})
                    </Menu.Item>
                  </>
                )}
              </Menu.Dropdown>
            </Menu>
          </Group>
        </Group>

        {/* Live Run Status */}
        {isRunning && run && (
          <LiveRunStatus run={run} onCancel={async () => {
            try {
              await api.cancelRun(run.id);
              setRun({ ...run, status: 'cancelled' });
              notifications.show({ title: 'Cancelled', message: 'Discovery cancelled', color: 'orange' });
            } catch (e: unknown) {
              notifications.show({ title: 'Error', message: (e as Error).message, color: 'red' });
            }
          }} />
        )}

        {/* Quick Stats */}
        {latestDiscovery && (
          <Grid>
            <Grid.Col span={{ base: 6, md: 3 }}>
              <Card withBorder p="md" ta="center">
                <Text size="xl" fw={700} c="blue">{discoveries.length}</Text>
                <Text size="sm" c="dimmed">Total Runs</Text>
              </Card>
            </Grid.Col>
            <Grid.Col span={{ base: 6, md: 3 }}>
              <Card withBorder p="md" ta="center">
                <Text size="xl" fw={700} c="violet">
                  {discoveries.reduce((sum, d) => sum + (d.summary?.total_insights || 0), 0)}
                </Text>
                <Text size="sm" c="dimmed">Total Insights</Text>
              </Card>
            </Grid.Col>
            <Grid.Col span={{ base: 6, md: 3 }}>
              <Card withBorder p="md" ta="center">
                <Text size="xl" fw={700} c="green">{latestDiscovery.summary?.total_insights || 0}</Text>
                <Text size="sm" c="dimmed">Latest Insights</Text>
              </Card>
            </Grid.Col>
            <Grid.Col span={{ base: 6, md: 3 }}>
              <Card withBorder p="md" ta="center">
                <Text size="xl" fw={700}>{latestDiscovery.total_steps}</Text>
                <Text size="sm" c="dimmed">Latest Steps</Text>
              </Card>
            </Grid.Col>
          </Grid>
        )}

        {/* Empty State */}
        {!latestDiscovery && !isRunning && (
          <Card withBorder p="xl" ta="center">
            <Stack align="center" gap="md">
              <IconSearch size={48} color="var(--mantine-color-gray-4)" />
              <Title order={3} c="dimmed">No discoveries yet</Title>
              <Text c="dimmed">Run your first discovery to start finding insights.</Text>
            </Stack>
          </Card>
        )}

        {/* Discovery History */}
        {discoveries.length > 0 && (
          <>
            <Title order={3}>Discoveries</Title>
            <Stack gap="sm">
              {discoveries.map((d) => (
                <Card key={d.id} withBorder p="md" radius="md" component={Link}
                  href={`/projects/${id}/discoveries/${d.id}`}
                  style={{ textDecoration: 'none', cursor: 'pointer' }}>
                  <Group justify="space-between">
                    <Group gap="sm">
                      <Text size="sm" fw={600}>
                        {new Date(d.discovery_date).toLocaleDateString('en-US', {
                          month: 'short', day: 'numeric', year: 'numeric',
                          hour: '2-digit', minute: '2-digit',
                        })}
                      </Text>
                      <Badge size="sm" variant="light"
                        color={d.run_type === 'partial' ? 'violet' : 'blue'}>
                        {d.run_type || 'full'}
                      </Badge>
                      {d.areas_requested && d.areas_requested.length > 0 && (
                        <Text size="xs" c="dimmed">{d.areas_requested.join(', ')}</Text>
                      )}
                    </Group>
                    <Group gap="sm">
                      <Badge size="sm" variant="outline" color="teal">
                        {d.summary?.total_insights || 0} insights
                      </Badge>
                      <Badge size="sm" variant="outline" color="gray">
                        {d.total_steps} steps
                      </Badge>
                    </Group>
                  </Group>
                </Card>
              ))}
            </Stack>
          </>
        )}
      </Stack>
    </Shell>
  );
}

function LiveRunStatus({ run, onCancel }: { run: DiscoveryRunStatus; onCancel: () => void }) {
  return (
    <Card withBorder p="lg" shadow="sm">
      <Group justify="space-between" mb="sm">
        <Group><Loader size="sm" /><Title order={4}>Discovery Running</Title></Group>
        <Group>
          <Badge color="blue" variant="light">{run.phase}</Badge>
          <Button size="xs" variant="light" color="red" onClick={onCancel}>Cancel</Button>
        </Group>
      </Group>
      <Progress value={run.progress} mb="sm" animated />
      <Text size="sm" c="dimmed" mb="md">{run.phase_detail}</Text>
      <Group gap="xl" mb="md">
        <Text size="xs" c="dimmed">Queries: {run.total_queries}</Text>
        <Text size="xs" c="dimmed">Insights: {run.insights_found}</Text>
      </Group>
      {run.steps && run.steps.length > 0 && (
        <ScrollArea h={200} type="auto">
          <Timeline active={run.steps.length - 1} bulletSize={18} lineWidth={2}>
            {run.steps.slice(-15).map((step, idx) => (
              <Timeline.Item key={idx}
                bullet={step.type === 'insight' ? <IconBulb size={10} /> :
                        step.type === 'error' ? <IconX size={10} /> :
                        <IconDatabase size={10} />}
                color={step.type === 'error' ? 'red' : step.type === 'insight' ? 'green' : 'blue'}
                title={<Text size="xs">{step.message}</Text>}>
                {step.llm_thinking && <Text size="xs" c="dimmed">{step.llm_thinking}</Text>}
              </Timeline.Item>
            ))}
          </Timeline>
        </ScrollArea>
      )}
    </Card>
  );
}
