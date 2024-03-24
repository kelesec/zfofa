import {Drawer} from "@arco-design/web-vue";
import {types} from "../../wailsjs/go/models";
import Asset = types.Asset;

/**
 * 表格中每一列的列类型
 */
export const COLUMN_DATA = [
    {
        title: 'ID', dataIndex: 'id', width: 70,
        sortable: {sortDirections: ['ascend', 'descend']}
    },
    {
        title: 'Host', dataIndex: 'host', ellipsis: true, width: 200, tooltip: true,
        sortable: {sortDirections: ['ascend', 'descend']}, slotName: 'host'
    },
    {
        title: 'Ip', dataIndex: 'ip', width: 170,
        sortable: {sortDirections: ['ascend', 'descend']}, slotName: 'ip'
    },
    {
        title: 'Port', dataIndex: 'port', width: 80,
        sortable: {sortDirections: ['ascend', 'descend']}
    },
    {
        title: 'Title', dataIndex: 'title',
        ellipsis: true, width: 200, tooltip: true
    },
    {
        title: 'Server', dataIndex: 'server',
        ellipsis: true, width: 120, tooltip: true
    },
    {
        title: 'Country', dataIndex: 'country',
        ellipsis: true, width: 200, tooltip: true
    },
    {
        title: 'Org', dataIndex: 'organization',
        ellipsis: true, width: 200, tooltip: true
    },
    {
        title: 'Header', dataIndex: 'header',
        ellipsis: true, width: 150, tooltip: true, slotName: 'header'
    },
    {
        title: 'Cert', dataIndex: 'certificate',
        ellipsis: true, width: 150, tooltip: true, slotName: 'cert'
    },
    {title: 'Alive', dataIndex: 'alive', width: 70},
    {
        title: 'LastUpdate', dataIndex: 'lastUpdateTime', ellipsis: true,
        width: 170, tooltip: true, sortable: {sortDirections: ['ascend', 'descend']}
    }
]


/**
 * 每条资产类型
 */
export class DataType extends Asset {
    id: number

    constructor(id: number, data: any) {
        super(data)
        this.id = id
    }
}


/**
 * 打开右侧抽屉展示信息
 * @param title 标题
 * @param content 内容
 * @constructor
 */
export function ShowDrawer(title: string, content: string) {
    Drawer.open({
        title: title,
        content: content,
        width: 350,
    })
}