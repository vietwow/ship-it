import React from 'react'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'
import ExpandIcon from '@material-ui/icons/ExpandMore'
import IconButton from '@material-ui/core/IconButton'
import SelectionDialog from './SelectionDialog'

import TimePassed from '../assets/time_passed.png'
import SlackIcon from '../assets/slack_icon.png'
import DockerIcon from '../assets/docker_icon.png'
import StatusChip from './StatusChip';
import SlackInfoDialog from './SlackInfoDialog';
import ImageInfoDialog from './ImageInfoDialog';

const imgAlt = "not found"

class SingleGridCell extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      expanded: false,
      selected_id: '',
      window_width: window.innerWidth,
      repoSelector: false,
      slackInfo: false,
      imageInfo: false,
    }
  }

  cellClick = (event) => {
    this.props.handleCellClick(event)
  }

  slackClicked = () => {
    this.setState({slackInfo: true})
  }

  slackClosed = () => {
    this.setState({slackInfo: false})
  }

  dockerClicked = () => {
    this.setState({imageInfo: true})
  }

  dockerClosed = () => {
    this.setState({imageInfo: false})
  }

  render() {
    const singleGridCellStyle = {
      backgroundSize: this.props.cellSize,
      width: this.props.cellSize,
      height: this.props.cellSize,
      display: 'inline-block',
      margin: this.props.cellMargin,
      marginBottom: 25,
      position: 'relative'
    }

    const cardStyle = {
      width: this.props.cellSize,
      height: this.props.cellSize
    }

    let deployDate = new Date(this.props.SingleGridCellData.lastDeployed)
    return (
      <div style={singleGridCellStyle} id={this.props.id} className='SingleGridCell'>
        <div>
          <Card style={cardStyle}>
            <CardContent>
              <Typography variant="h5" component="h2">
                {this.props.SingleGridCellData.name}
              </Typography>
              <StatusChip status={this.props.SingleGridCellData.deployment.status} />
              <div>
                <IconButton>
                  <img src={TimePassed} alt={imgAlt} />
                </IconButton>
                {deployDate.toDateString()}
              </div>
              <div className='row-align'>
                <SelectionDialog source={this.props.SingleGridCellData.code.github} chart={this.props.SingleGridCellData.artifacts.chart.path}/>
                <IconButton onClick={this.slackClicked}>
                  <img src={SlackIcon} width="32" height="32" alt={imgAlt} />
                </IconButton>
                {
                  this.state.slackInfo ? <SlackInfoDialog open={this.state.slackInfo} owner={this.props.SingleGridCellData.owner.team} slack={this.props.SingleGridCellData.owner.slack} handleClose={this.slackClosed}/> : null
                }
                <IconButton onClick={this.dockerClicked}>
                  <img src={DockerIcon} width="32" height="32" alt={imgAlt} />
                </IconButton>
                {
                  this.state.imageInfo ? <ImageInfoDialog open={this.state.imageInfo} handleClose={this.dockerClosed} docker={this.props.SingleGridCellData.artifacts.docker} /> : null
                }
              </div>
              <div>
                <IconButton onClick={this.cellClick}>
                  <ExpandIcon />
                </IconButton>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    )
  }
}

export default SingleGridCell