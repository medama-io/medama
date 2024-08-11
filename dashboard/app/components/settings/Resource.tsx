import * as Progress from '@radix-ui/react-progress';
import byteSize from 'byte-size';

import { Group } from '@/components/layout/Flex';

import classes from './Resource.module.css';

interface PanelProps {
	title: string;
	usage: number;
}

interface ResourcePanelProps {
	title: string;
	used: number;
	total: number;
}

interface ResourcePanelCPUProps extends PanelProps {
	cores: number;
	threads: number;
}

const Panel = ({
	title,
	usage,
	children,
}: React.PropsWithChildren<PanelProps>) => (
	<div className={classes.panel}>
		<Group style={{ marginBottom: 16 }}>
			<p style={{ fontSize: 16, fontWeight: 500 }}>{title}</p>
			<p style={{ fontSize: 22, fontFamily: 'monospace' }}>
				{usage.toFixed(2)}%
			</p>
		</Group>
		<Progress.Root className={classes.progress} value={Math.round(usage)}>
			<Progress.Indicator
				className={classes.indicator}
				style={{ transform: `translateX(-${100 - Math.round(usage)}%)` }}
			/>
		</Progress.Root>
		{children}
	</div>
);

export const ResourcePanel = ({ title, used, total }: ResourcePanelProps) => {
	const usage = (used / total) * 100;
	return (
		<Panel title={title} usage={usage}>
			<div className={classes.group}>
				<p>
					Used: <span>{String(byteSize(used))}</span>
				</p>
				<p>
					Capacity: <span>{String(byteSize(total))}</span>
				</p>
			</div>
		</Panel>
	);
};

export const ResourcePanelCPU = ({
	title,
	usage,
	cores,
	threads,
}: ResourcePanelCPUProps) => (
	<Panel title={title} usage={usage}>
		<div className={classes.group}>
			<p>
				Cores: <span>{String(cores)}</span>
			</p>
			<p>
				Threads: <span>{String(threads)}</span>
			</p>
		</div>
	</Panel>
);
