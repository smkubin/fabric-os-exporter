## porterrshow metrics

| # | command | Metrics Name | Labels | Description |
| -- | -- | --| --| --| 
| ~~01~~ | ~~porterrshow~~ | ~~fabricos_portstats_frames_tx~~| ~~resource,portIndex~~ | ~~Number of frames transmitted (Tx) errors.~~|
| ~~02~~ | ~~porterrshow~~ |~~fabricos_portstats_frames_rx~~  | ~~resource,portIndex~~| ~~Number of frames received (Rx) errors.~~ |
| ~~03~~ | ~~porterrshow~~| ~~fabricos_portstats_enc_in~~ | ~~resource,portIndex~~ | ~~Number of encoding errors inside frames received (Rx).~~ |
| 04 | porterrshow | fabricos_portstats_crc_err | resource,portIndex | Number of frames with CRC errors received (Rx). |
| 05 | porterrshow | fabricos_portstats_crc_g_eof| resource,portIndex | Number of frames with CRC errors with good EOF received (Rx). |
| ~~06~~ | ~~porterrshow~~ |~~fabricos_portstats_too_short~~ | ~~resource,portIndex~~ | ~~Number of frames shorter than minimum received (Rx).~~ |
| ~~07~~ | ~~porterrshow~~ |~~fabricos_portstats_too_long~~ | ~~resource,portIndex~~ | ~~Number of frames longer than maximum received (Rx).~~ |
| ~~08~~ | ~~porterrshow~~ |~~fabricos_portstats_bad_eof~~ | ~~resource,portIndex~~ | ~~Number of frames with bad end-of-frame delimiters received (Rx).~~ | 
| 09 | porterrshow |fabricos_portstats_enc_out | resource,portIndex | Number of encoding error outside of frames received (Rx). |
| ~~10~~ | ~~porterrshow~~ |~~fabricos_portstats_disc_c3~~ | ~~resource,portIndex~~ | ~~Number of Class 3 frames discarded (Rx).~~ | 
| ~~11~~ | ~~porterrshow~~ |~~fabricos_portstats_link_fail~~ | ~~resource,portIndex~~ | ~~Number of link failures (LF1 or LF2 states) received (Rx).~~ |
| ~~12~~ | ~~porterrshow~~ |~~fabricos_portstats_loss_sync~~ | ~~resource,portIndex~~ | ~~Number of times synchronization was lost (Rx).~~ |
| ~~13~~ | ~~porterrshow~~ |~~fabricos_portstats_loss_sig~~ | ~~resource,portIndex~~ | ~~Number of times a loss of signal was received (increments whenever an SFP is removed) (Rx).~~ |
| ~~14~~ | ~~porterrshow~~ |~~fabricos_portstats_frjt~~ | ~~resource,portIndex~~ | ~~Number of transmitted frames rejected with F_RJT (Tx).~~|  
| ~~15~~ | ~~porterrshow~~ |~~fabricos_portstats_fbsy~~ | ~~resource,portIndex~~ | ~~Number of transmitted frames busied with F_BSY (Tx).~~ |
| ~~16~~ | ~~porterrshow~~ |~~fabricos_portstats_c3_timeout_tx~~ | ~~resource,portIndex~~ | ~~The number of transmit class 3 frames discarded at the transmission port due to timeout (platform- and port-specific).~~ |
| ~~17~~ | ~~porterrshow~~ |~~fabricos_portstats_c3_timeout_rx~~ | ~~resource,portIndex~~ | ~~The number of receive class 3 frames received at this port and discarded at the transmission port due to timeout (platform- and port-specific).~~ |
| 18 | porterrshow |fabricos_portstats_pcs_err | resource,portIndex | The number of Physical Coding Sublayer (PCS) block errors. This counter records encoding violations on 10 Gbps or 16 Gbps ports. |
| 19 | porterrshow |fabricos_portstats_uncor_err_fec | resource,portIndex | The number of uncorrectable forward error corrections (FEC). |
| 20 | portstatsshow | fabricos_portstats_cor_fec | resource,portIndex | Count of blocks that were corrected by FEC. |