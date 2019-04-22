## porterrshow metrics
| # | command | Metrics Name | Labels | Description |
| -- | -- | --| --| --| 
| 01 | porterrshow | fabricos_porterr_frames_tx| target,portNo | Number of frames transmitted (Tx).|
| 02 | porterrshow |fabricos_porterr_frames_rx  | target,portNo| Number of frames received (Rx). |
| 03 | porterrshow| fabricos_porterr_enc_in | target,portNo | Number of encoding errors inside frames received (Rx). |
| 04 | porterrshow | fabricos_porterr_crc_err | target,portNo | Number of frames with CRC errors received (Rx). |
| 05 | porterrshow | fabricos_porterr_crc_g_eof| target,portNo | Number of frames with CRC errors with good EOF received (Rx). |
| 06 | porterrshow |fabricos_porterr_too_short | target,portNo | Number of frames shorter than minimum received (Rx). |
| 07 | porterrshow |fabricos_porterr_too_long | target,portNo | Number of frames longer than maximum received (Rx). |
| 08 | porterrshow |fabricos_porterr_bad_eof | target,portNo | Number of frames with bad end-of-frame delimiters received (Rx). | 
| 09 | porterrshow |fabricos_porterr_enc_out | target,portNo | Number of encoding error outside of frames received (Rx). |
| 10 | porterrshow |fabricos_porterr_disc_c3 | target,portNo | Number of Class 3 frames discarded (Rx). | 
| 11 | porterrshow |fabricos_porterr_link_fail | target,portNo | Number of link failures (LF1 or LF2 states) received (Rx). |
| 12 | porterrshow |fabricos_porterr_loss_sync | target,portNo | Number of times synchronization was lost (Rx). |
| 13 | porterrshow |fabricos_porterr_loss_sig | target,portNo | Number of times a loss of signal was received (increments whenever an SFP is removed) (Rx). |
| 14 | porterrshow |fabricos_porterr_frjt | target,portNo | Number of transmitted frames rejected with F_RJT (Tx).|  
| 15 | porterrshow |fabricos_porterr_fbsy | target,portNo | Number of transmitted frames busied with F_BSY (Tx). |
| 16 | porterrshow |fabricos_porterr_c3_timeout_tx | target,portNo | The number of transmit class 3 frames discarded at the transmission port due to timeout (platform- and port-specific). |
| 17 | porterrshow |fabricos_porterr_c3_timeout_rx | target,portNo | The number of receive class 3 frames received at this port and discarded at the transmission port due to timeout (platform- and port-specific). |
| 18 | porterrshow |fabricos_porterr_pcs_err | target,portNo | The number of Physical Coding Sublayer (PCS) block errors. This counter records encoding violations on 10 Gbps or 16 Gbps ports. |
| 19 | porterrshow |fabricos_porterr_uncor_err | target,portNo | The number of uncorrectable forward error corrections (FEC). |